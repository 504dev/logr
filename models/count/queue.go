package count

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/clickhouse"
	"github.com/504dev/logr/clickhouse/queue"
	"time"
)

var Queue *queue.Queue

func RunQueue() {
	sql := `
		INSERT INTO counts (day, timestamp, dash_id, hostname, logname, keyname, version, inc, max, min, avg_sum, avg_num, per_tkn, per_ttl, time_dur)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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

func PushToQueue(c *_types.Count) error {
	Queue.Push(c.AsVector())
	return nil
}
