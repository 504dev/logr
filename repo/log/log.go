package log

import (
	"fmt"
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/dbs/clickhouse"
	"github.com/504dev/logr/dbs/clickhouse/batcher"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"github.com/jmoiron/sqlx"
	"time"
)

type LogRepo struct {
	conn    *sqlx.DB
	batcher *batcher.Batcher[*_types.Log]
}

func NewLogRepo() *LogRepo {
	return &LogRepo{
		conn: clickhouse.Conn(),
	}
}

func (repo *LogRepo) Insert(log *_types.Log) error {
	return repo.BatchInsert([]*_types.Log{log})
}

func (repo *LogRepo) BatchInsert(batch []*_types.Log) error {
	const sql = `
		INSERT INTO logs (day, timestamp, dash_id, hostname, logname, level, message, pid, version)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	tx, err := repo.conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range batch {
		day := time.Unix(0, v.Timestamp).UTC().Format("2006-01-02")
		values := []any{day, v.Timestamp, v.DashId, v.Hostname, v.Logname, v.Level, v.Message, v.Pid, v.Version}
		_, err = stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (repo *LogRepo) GetByFilter(filter types.Filter) (_types.Logs, error) {
	where, values, err := filter.ToSQL()
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
		whereToday := fmt.Sprintf("%v AND day = '%v'", where, day)
		logsToday, err := repo.getByFilter(whereToday, values, limit)
		if err != nil {
			return nil, err
		}

		rest := limit - len(logsToday)
		if rest == 0 {
			return logsToday, nil
		}
		whereRest := fmt.Sprintf("%v AND day < '%v'", where, day)
		logsRest, err := repo.getByFilter(whereRest, values, rest)
		if err != nil {
			return nil, err
		}

		return append(logsToday, logsRest...), nil
	}

	return repo.getByFilter(where, values, limit)
}

func (repo *LogRepo) getByFilter(where string, values []interface{}, limit int) (_types.Logs, error) {
	delta := Logger.Duration()
	sql := `
      SELECT timestamp, dash_id, hostname, logname, level, message
      FROM logs ` + where + `
      ORDER BY day DESC, timestamp DESC
      LIMIT ` + fmt.Sprint(limit)
	logs := _types.Logs{}
	err := repo.conn.Select(&logs, sql, values...)
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
	err := repo.conn.Select(&stats, sql, dashId, logname)
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
	err := repo.conn.Select(&stats, sql, dashId)
	if err != nil {
		return nil, err
	}
	Logger.Debug("%v %v", sql, dashId)
	return stats, nil
}
