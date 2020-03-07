package log

import (
	"fmt"
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/types"
	"time"
)

func Create(log *types.Log) error {
	conn := clickhouse.Conn()

	day := time.Unix(0, log.Timestamp).Format("2006-01-02")
	values := []interface{}{day, log.Timestamp, log.DashId, log.Hostname, log.Logname, log.Level, log.Message}
	sql := `INSERT INTO logs (day, timestamp, dash_id, hostname, logname, level, message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

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

	logger.Debug("%v %v", sql, values)
	logs := types.Logs{}
	err := conn.Select(&logs, sql, values...)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func GetDashStats(dashIds []int) ([]*types.DashStatRow, error) {
	conn := clickhouse.Conn()
	sql := `
      SELECT dash_id, hostname, logname, level, count(*) AS cnt, max(day) AS updated
      FROM logs WHERE dash_id IN (?)
      GROUP BY dash_id, hostname, logname, level
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
