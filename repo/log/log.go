package log

import (
	"fmt"
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/dbs/clickhouse"
	"github.com/504dev/logr/dbs/clickhouse/queue"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"time"
)

func NewLogRepo() *LogRepo {
	sql := `
		INSERT INTO logs (day, timestamp, dash_id, hostname, logname, level, message, pid, version)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return &LogRepo{
		queue: queue.NewQueue(&queue.QueueConfig{
			DB:            clickhouse.Conn(),
			Sql:           sql,
			FlushInterval: time.Second,
			FlushCount:    1000,
		}),
	}
}

func (repo *LogRepo) GetByFilter(filter types.Filter) (_types.Logs, error) {
	where, values, err := filter.ToSql()
	if err != nil {
		return nil, err
	}
	limit := filter.Limit
	if limit == 0 {
		limit = 100
	}
	if filter.Pattern == "" {
		now := time.Now()
		if filter.Offset != 0 {
			now = time.Unix(0, filter.Offset)
		}
		day := now.UTC().Format(time.RFC3339)[0:10]
		w1 := fmt.Sprintf("%v AND day = '%v'", where, day)
		logs, err := repo.getByFilter(w1, values, limit)
		if err != nil {
			return nil, err
		}
		rest := limit - len(logs)
		if rest == 0 {
			return logs, nil
		}
		w2 := fmt.Sprintf("%v AND day < '%v'", where, day)
		tmp, err := repo.getByFilter(w2, values, rest)
		if err != nil {
			return nil, err
		}
		return append(logs, tmp...), nil
	}
	return repo.getByFilter(where, values, limit)
}

func (repo *LogRepo) getByFilter(where string, values []interface{}, limit int) (_types.Logs, error) {
	conn := clickhouse.Conn()
	delta := Logger.Duration()
	sql := `
      SELECT timestamp, dash_id, hostname, logname, level, message
      FROM logs ` + where + `
      ORDER BY day DESC, timestamp DESC
      LIMIT ` + fmt.Sprint(limit)
	logs := _types.Logs{}
	err := conn.Select(&logs, sql, values...)
	Logger.Debug("%v\n%v\ncount: %v, time: %v", sql, values, len(logs), delta())
	return logs, err
}

func (repo *LogRepo) StatsByLogname(dashId int, logname string) ([]*types.DashStatRow, error) {
	sql := `
      SELECT hostname, level, version, count(*) AS cnt, max(timestamp) AS updated
      FROM logs
      WHERE dash_id = ? AND logname = ? AND day > toDate(now() - interval 7 day)
      GROUP BY hostname, level, version
    `
	stats := types.DashStatRows{}
	err := clickhouse.Conn().Select(&stats, sql, dashId, logname)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (repo *LogRepo) StatsByDashboard(dashId int) ([]*types.DashStatRow, error) {
	sql := `
      SELECT logname, level, count(*) AS cnt, max(timestamp) AS updated, substring(anyLast(message), 1, 100) as message
      FROM logs
      WHERE
        dash_id = ?
        AND day >= toDate(now() - interval 1 day)
        AND timestamp > toUnixTimestamp(now() - interval 1 day) * 1e9
      GROUP BY logname, level
    `
	stats := types.DashStatRows{}
	err := clickhouse.Conn().Select(&stats, sql, dashId)
	if err != nil {
		return nil, err
	}
	Logger.Debug("%v %v", sql, dashId)
	return stats, nil
}
