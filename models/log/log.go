package log

import "time"

var logs = Logs{
	{
		Timestamp: time.Now(),
		DashId:    1,
		Host:      "localhost",
		Project:   "kodlog.project",
		Logname:   "performance.log",
		Level:     1,
		Message:   "504",
		Order:     0,
	},
}

func GetLast() *Log {
	return logs[len(logs)-1]
}
