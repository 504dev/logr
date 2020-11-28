package log

import (
	"fmt"
	"github.com/504dev/logr/clickhouse"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"time"
)

func GetByFilter(f types.Filter) (types.Logs, error) {
	where, values, err := f.ToSql()
	if err != nil {
		return nil, err
	}
	limit := f.Limit
	if limit == 0 {
		limit = 100
	}
	if f.Pattern == "" && f.Offset == 0 {
		today := time.Now().UTC().Format(time.RFC3339)[0:10]
		w1 := fmt.Sprintf("%v AND day = '%v'", where, today)
		logs, err := getByFilter(w1, values, limit)
		if err != nil {
			return nil, err
		}
		rest := limit - len(logs)
		if rest == 0 {
			return logs, nil
		}
		w2 := fmt.Sprintf("%v AND day < '%v'", where, today)
		tmp, err := getByFilter(w2, values, rest)
		if err != nil {
			return nil, err
		}
		return append(logs, tmp...), nil
	}
	return getByFilter(where, values, limit)
}

func getByFilter(where string, values []interface{}, limit int) (types.Logs, error) {
	conn := clickhouse.Conn()
	sql := `
      SELECT timestamp, dash_id, hostname, logname, level, message
      FROM logs ` + where + `
      ORDER BY day DESC, timestamp DESC
      LIMIT ` + fmt.Sprint(limit)
	logs := types.Logs{}
	err := conn.Select(&logs, sql, values...)
	Logger.Debug("%v\n%v\n%v", sql, values, len(logs))
	return logs, err
}

func GetDashStats(dashId int) ([]*types.DashStatRow, error) {
	sql := `
      SELECT hostname, logname, level, version, count(*) AS cnt, max(timestamp) AS updated
      FROM logs
      WHERE dash_id = ? AND day > toDate(now() - interval 7 day)
      GROUP BY hostname, logname, level, version
    `
	stats := types.DashStatRows{}
	err := clickhouse.Conn().Select(&stats, sql, dashId)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func GetDashLognames(dashId int) ([]*types.DashStatRow, error) {
	sql := `
      SELECT logname, count(*) AS cnt
      FROM logs
      WHERE
        dash_id = ?
        AND day >= toDate(now() - interval 1 day)
        AND timestamp > toUnixTimestamp(now() - interval 1 hour) * 1e9
      GROUP BY logname
    `
	stats := types.DashStatRows{}
	err := clickhouse.Conn().Select(&stats, sql, dashId)
	if err != nil {
		return nil, err
	}
	Logger.Debug("%v %v", sql, dashId)
	return stats, nil
}
