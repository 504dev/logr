package log

import (
	"fmt"
	"github.com/504dev/kidlog/clickhouse"
)

func GetAll(f Filter) (Logs, error) {
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
	logs := make(Logs, 0)

	for rows.Next() {
		var log Log
		err := rows.StructScan(&log)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}
	return logs, nil
}

type DashStatRow struct {
	Hostname string `db:"hostname"  json:"hostname"`
	Logname  string `db:"logname"   json:"logname"`
	Level    string `db:"level"     json:"level"`
	Cnt      int    `db:"cnt"       json:"cnt"`
	Updated  string `db:"updated"   json:"updated"`
}

func GetDashStats(dashId int) ([]*DashStatRow, error) {
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
	stats := make([]*DashStatRow, 0)

	for rows.Next() {
		var row DashStatRow
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &row)
	}
	return stats, nil
}
