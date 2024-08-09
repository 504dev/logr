package count

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/clickhouse"
	"github.com/504dev/logr/clickhouse/queue"
	"github.com/504dev/logr/types"
	"time"
)

type CountStorage struct {
	queue *queue.Queue
}

func NewCountStorage() *CountStorage {
	sql := `
		INSERT INTO counts (day, timestamp, dash_id, hostname, logname, keyname, version, inc, max, min, avg_sum, avg_num, per_tkn, per_ttl, time_dur)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return &CountStorage{
		queue: queue.NewQueue(&queue.QueueConfig{
			DB:            clickhouse.Conn(),
			Sql:           sql,
			FlushInterval: time.Second,
			FlushCount:    1000,
		}),
	}
}

func (storage *CountStorage) RunQueue() *CountStorage {
	storage.queue.Run()
	return storage
}

func (storage *CountStorage) StopQueue() error {
	return storage.queue.Stop()
}

func (storage *CountStorage) Store(c *_types.Count) error {
	cv := (*types.CountVector)(c)
	storage.queue.Push(cv.AsVector())
	return nil
}
