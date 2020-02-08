package log

import (
	"fmt"
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/types"
	"time"
)

func Create(log *types.Log) error {
	day := time.Unix(0, log.Timestamp).Format("2006-01-02")
	values := []interface{}{day, log.Timestamp, log.DashId, log.Hostname, log.Logname, log.Level, log.Message}
	conn := clickhouse.Conn()

	sqlstr := `INSERT INTO logs (day, timestamp, dash_id, hostname, logname, level, message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	tx, _ := conn.Begin()
	stmt, _ := tx.Prepare(sqlstr)
	defer stmt.Close()

	_, err := stmt.Exec(values...)
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
	sql := `
      SELECT timestamp, dash_id, hostname, logname, level, message
      FROM logs ` + where + `
      ORDER BY day DESC, timestamp DESC
    `
	if f.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT %v", f.Limit)
	}
	fmt.Println(sql, values)
	rows, err := conn.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	logs := make(types.Logs, 0)

	for rows.Next() {
		var log types.Log
		err := rows.StructScan(&log)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}
	return logs, nil
}

func GetDashStats(dashId int) ([]*types.DashStatRow, error) {
	conn := clickhouse.Conn()
	sql := `
      SELECT hostname, logname, level, count(*) AS cnt, max(day) AS updated
      FROM logs WHERE dash_id = ?
      GROUP BY hostname, logname, level
    `
	rows, err := conn.Queryx(sql, dashId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stats := make([]*types.DashStatRow, 0)

	for rows.Next() {
		var row types.DashStatRow
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &row)
	}
	return stats, nil
}
