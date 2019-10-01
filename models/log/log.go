package log

import "time"

var logs = Logs{
	{
		DashId:    1,
		Timestamp: time.Now(),
		Host:      "localhost",
		Project:   "kodlog.project",
		Logname:   "performance.log",
		Level:     1,
		Message:   "Hello",
	},
}

func GetAll() Logs {
	return logs
}

func GetLast() *Log {
	return logs[len(logs)-1]
}
