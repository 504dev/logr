package log

import (
	"github.com/504dev/kidlog/clickhouse"
	"time"
)

var logs = Logs{
	{
		DashId:    1,
		Timestamp: time.Now().UnixNano(),
		Hostname:  "localhost",
		Logname:   "performance.log",
		Level:     0,
		Message:   "Hello",
	},
}

func GetAll(f Filter) (Logs, error) {
	conn := clickhouse.Conn()
	where, values := f.ToSql()
	sql := "SELECT timestamp, dash_id, hostname, logname, level, message FROM logs " + where
	rows, err := conn.Queryx(sql, values)
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

func GetLast() *Log {
	return logs[len(logs)-1]
}
