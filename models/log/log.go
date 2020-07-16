package log

import (
	"fmt"
	"github.com/504dev/logr/clickhouse"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
)

func GetByFilter(f types.Filter) (types.Logs, error) {
	conn := clickhouse.Conn()
	where, values := f.ToSql()
	limit := f.Limit
	if limit == 0 {
		limit = 100
	}
	sql := `
      SELECT timestamp, dash_id, hostname, logname, level, message
      FROM logs ` + where + `
      ORDER BY day DESC, timestamp DESC
      LIMIT ` + fmt.Sprint(limit)

	Logger.Debug("%v %v", sql, values)
	logs := types.Logs{}
	err := conn.Select(&logs, sql, values...)
	if err != nil {
		return nil, err
	}
	return logs, nil
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
      SELECT
        logname, count(*) AS cnt FROM logs
      WHERE
        dash_id = ?
        AND day >= toDate(now() - interval 1 day)
        AND timestamp > toUnixTimestamp(now() - interval 1 hour) * 1e9
      GROUP BY
        logname
    `
	stats := types.DashStatRows{}
	err := clickhouse.Conn().Select(&stats, sql, dashId)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
