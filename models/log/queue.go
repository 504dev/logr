package log

import (
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/clickhouse/queue"
	"github.com/504dev/kidlog/types"
	"time"
)

var Queue *queue.Queue

func RunQueue() {
	Queue = queue.NewQueue(&queue.Config{
		DB:            clickhouse.Conn(),
		Sql:           "INSERT INTO logs (day, timestamp, dash_id, hostname, logname, level, message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		FlushInterval: time.Second,
		FlushCount:    1000,
	})
	Queue.Run()
}

func StopQueue() error {
	return Queue.Stop()
}

func PushToQueue(log *types.Log) error {
	day := time.Unix(0, log.Timestamp).Format("2006-01-02")
	values := []interface{}{day, log.Timestamp, log.DashId, log.Hostname, log.Logname, log.Level, log.Message}
	Queue.Push(values)
	return nil
}
