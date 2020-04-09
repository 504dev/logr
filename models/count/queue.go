package count

import (
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/clickhouse/queue"
	"github.com/504dev/kidlog/types"
	"time"
)

var Queue *queue.Queue

func RunQueue() {
	sql := `
		INSERT INTO counts (day, timestamp, dash_id, hostname, logname, keyname, version, inc, max, min, avg_sum, avg_num, per_tkn, per_ttl)
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

func PushToQueue(c *types.Count) error {
	Queue.Push(c.AsVector())
	return nil
}
