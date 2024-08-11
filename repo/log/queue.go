package log

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/dbs/clickhouse"
	"github.com/504dev/logr/dbs/clickhouse/queue"
	"time"
)

type LogStorage struct {
	queue *queue.Queue
}

func NewLogStorage() *LogStorage {
	sql := `
		INSERT INTO logs (day, timestamp, dash_id, hostname, logname, level, message, pid, version)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return &LogStorage{
		queue: queue.NewQueue(&queue.QueueConfig{
			DB:            clickhouse.Conn(),
			Sql:           sql,
			FlushInterval: time.Second,
			FlushCount:    1000,
		}),
	}
}

func (storage *LogStorage) RunQueue() *LogStorage {
	storage.queue.Run()
	return storage
}

func (storage *LogStorage) StopQueue() error {
	return storage.queue.Stop()
}

func (storage *LogStorage) Store(log *_types.Log) error {
	day := time.Unix(0, log.Timestamp).UTC().Format("2006-01-02")
	values := []interface{}{day, log.Timestamp, log.DashId, log.Hostname, log.Logname, log.Level, log.Message, log.Pid, log.Version}
	storage.queue.Push(values)
	return nil
}
