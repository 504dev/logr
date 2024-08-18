package queue

import (
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
	config QueueConfig[T]
	mu     sync.Mutex
	batch  []T
	stop   chan struct{}
	done   chan struct{}
}

func NewQueue[T any](nlimit int, tlimit time.Duration, handler func([]T)) *Queue[T] {
	return &Queue[T]{
		config: QueueConfig[T]{
			Handler:       handler,
			FlushInterval: tlimit,
			FlushCount:    nlimit,
		},
		batch: make([]T, 0, nlimit),
		stop:  make(chan struct{}),
		done:  make(chan struct{}),
	}
}

func (q *Queue[T]) Run() {
	for {
		select {
		case <-time.After(q.config.FlushInterval):
			q.flushSafe()
		case <-q.stop:
			q.flushSafe()
			close(q.done)
			return
		}
	}
}

func (q *Queue[T]) Push(item T) bool {
	select {
	case <-q.stop:
		return false
	default:
		//
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	q.batch = append(q.batch, item)
	if len(q.batch) >= q.config.FlushCount {
		q.flush()
	}
	return true
}

func (q *Queue[T]) flush() int {
	if len(q.batch) == 0 {
		return 0
	}

	batch := q.batch
	q.batch = make([]T, 0, q.config.FlushCount)

	q.config.Handler(batch)

	return len(batch)
}

func (q *Queue[T]) flushSafe() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.flush()
}

func (q *Queue[T]) Stop() {
	close(q.stop)
	<-q.done
}

func (b *Queue[T]) HasStopped() <-chan struct{} {
	return b.stop
}

func (b *Queue[T]) Done() <-chan struct{} {
	return b.done
}
