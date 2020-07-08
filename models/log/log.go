package log

import (
	"fmt"
	"github.com/504dev/logr/clickhouse"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"time"
)

func GetByFilter(f types.Filter) (types.Logs, error) {
	duration := Logger.Time("response:/logs", time.Millisecond)
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
	duration()
	Logger.Inc("/logs:cnt", 1)
	return logs, nil
}

func GetDashStats(dashId int) ([]*types.DashStatRow, error) {
	conn := clickhouse.Conn()
	sql := `
      SELECT hostname, logname, level, version, count(*) AS cnt, max(timestamp) AS updated
      FROM logs
      WHERE dash_id = ? AND day > toDate(now() - interval 7 day)
      GROUP BY hostname, logname, level, version
    `
	stats := types.DashStatRows{}
	err := conn.Select(&stats, sql, dashId)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
