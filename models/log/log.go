package log

import (
	"fmt"
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/clickhouse"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"time"
)

func GetByFilter(f types.Filter) (_types.Logs, error) {
	where, values, err := f.ToSql()
	if err != nil {
		return nil, err
	}
	limit := f.Limit
	if limit == 0 {
		limit = 100
	}
	if f.Pattern == "" {
		now := time.Now()
		if f.Offset != 0 {
			now = time.Unix(0, f.Offset)
		}
		day := now.UTC().Format(time.RFC3339)[0:10]
		w1 := fmt.Sprintf("%v AND day = '%v'", where, day)
		logs, err := getByFilter(w1, values, limit)
		if err != nil {
			return nil, err
		}
		rest := limit - len(logs)
		if rest == 0 {
			return logs, nil
		}
		w2 := fmt.Sprintf("%v AND day < '%v'", where, day)
		tmp, err := getByFilter(w2, values, rest)
		if err != nil {
			return nil, err
		}
		return append(logs, tmp...), nil
	}
	return getByFilter(where, values, limit)
}

func getByFilter(where string, values []interface{}, limit int) (_types.Logs, error) {
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

func StatsByLogname(dashId int, logname string) ([]*types.DashStatRow, error) {
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

func StatsByDashboard(dashId int) ([]*types.DashStatRow, error) {
	sql := `
      SELECT logname, level, count(*) AS cnt, max(timestamp) AS updated, substring(anyLast(message), 1, 100) as message
      FROM logs
      WHERE
        dash_id = ?
        AND day >= toDate(now() - interval 1 day)
        AND timestamp > toUnixTimestamp(now() - interval 1 hour) * 1e9
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
