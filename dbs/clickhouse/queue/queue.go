package queue

import (
	. "github.com/504dev/logr/logger"
	"sync"
	"time"
)

// deprecated use clickhouse/batcher

type QueueConfig[T any] struct {
	Handler       func([]T)
	FlushInterval time.Duration
	FlushCount    int
}

type Queue[T any] struct {
	sync.Mutex
	*QueueConfig[T]
	batch []T
	stop  chan struct{}
}

func NewQueue[T any](c *QueueConfig[T]) *Queue[T] {
	return &Queue[T]{
		QueueConfig: c,
		batch:       make([]T, 0, c.FlushCount),
		stop:        make(chan struct{}),
	}
}

func (q *Queue[T]) Run() {
	go (func() {
		for {
			select {
			case <-time.After(q.FlushInterval):
				if err := q.Flush(); err != nil {
					Logger.Error(err)
				}
			case <-q.stop:
				return
			}
		}
	})()
}

func (q *Queue[T]) Stop() error {
	close(q.stop)
	return q.Flush()
}

func (q *Queue[T]) Push(item T) {
	q.Lock()
	q.batch = append(q.batch, item)
	count := len(q.batch)
	q.Unlock()
	if count >= q.FlushCount {
		q.Flush()
	}
}

func (q *Queue[T]) Flush() error {
	q.Lock()
	if len(q.batch) == 0 {
		q.Unlock()
		return nil
	}
	batch := q.batch
	q.batch = make([]T, 0, q.FlushCount)
	q.Unlock()

	q.Handler(batch)

	return nil
}
