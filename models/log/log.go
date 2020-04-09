package log

import (
	"fmt"
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/types"
	"time"
)

func GetByFilter(f types.Filter) (types.Logs, error) {
	ts := time.Now()
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

	logger.Debug("%v %v", sql, values)
	logs := types.Logs{}
	err := conn.Select(&logs, sql, values...)
	if err != nil {
		return nil, err
	}
	delta := time.Now().Sub(ts).Seconds()
	logger.Avg("/logs:time", delta).Max(delta).Min(delta)
	logger.Inc("/logs:cnt", 1)
	return logs, nil
}

func GetDashStats(dashIds []int) ([]*types.DashStatRow, error) {
	conn := clickhouse.Conn()
	sql := `
      SELECT dash_id, hostname, logname, level, version, count(*) AS cnt, max(day) AS updated
      FROM logs WHERE dash_id IN (?)
      GROUP BY dash_id, hostname, logname, level, version
    `
	stats := types.DashStatRows{}
	err := conn.Select(&stats, sql, dashIds)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func GetFrequentDashboards(len int) ([]int, error) {
	conn := clickhouse.Conn()
	sql := `
      SELECT dash_id, count(*) AS cnt
      FROM logs WHERE timestamp > ?
      GROUP BY dash_id
      ORDER BY cnt DESC
      LIMIT ?
    `
	ts := time.Now().Add(-3 * time.Hour)
	rows, err := conn.Queryx(sql, ts, len)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]int, 0, len)

	for rows.Next() {
		var id, cnt int
		err := rows.Scan(&id, &cnt)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}
