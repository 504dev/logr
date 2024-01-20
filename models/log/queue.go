package log

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/clickhouse"
	"github.com/504dev/logr/clickhouse/queue"
	"time"
)

var Queue *queue.Queue

func RunQueue() {
	sql := `
        INSERT INTO logs (day, timestamp, dash_id, hostname, logname, level, message, pid, version)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	Queue = queue.NewQueue(&queue.Config{
		DB:            clickhouse.Conn(),
		Sql:           sql,
		FlushInterval: time.Second,
		FlushCount:    1000,
	})
	Queue.Run()
}

func StopQueue() error {
	return Queue.Stop()
}

func PushToQueue(log *_types.Log) error {
	day := time.Unix(0, log.Timestamp).UTC().Format("2006-01-02")
	values := []interface{}{day, log.Timestamp, log.DashId, log.Hostname, log.Logname, log.Level, log.Message, log.Pid, log.Version}
	Queue.Push(values)
	return nil
}
