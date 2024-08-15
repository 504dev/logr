package batcher

import (
	"time"
)

type Batcher[T any] struct {
	nlimit  int
	tlimit  time.Duration
	channel chan T
	batch   []T
	handler func([]T)
	stop    chan struct{}
	done    chan struct{}
}

func NewBatcher[T any](nlimit int, tlimit time.Duration, handler func([]T)) *Batcher[T] {
	return &Batcher[T]{
		nlimit:  nlimit,
		tlimit:  tlimit,
		channel: make(chan T, nlimit),
		batch:   make([]T, 0, nlimit),
		handler: handler,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

func (b *Batcher[T]) Push(item T) bool {
	select {
	case <-b.stop:
		return false
	default:
		b.channel <- item
		return true
	}
}

func (b *Batcher[T]) Run() {
	ticker := time.NewTicker(b.tlimit)
	defer ticker.Stop()
	defer close(b.done)

	for {
		select {
		case item := <-b.channel:
			b.append(item)
		case <-ticker.C:
			b.flush()
		case <-b.stop:
			b.drain()
			return
		}
	}
}

func (b *Batcher[T]) append(item T) {
	b.batch = append(b.batch, item)
	if len(b.batch) == b.nlimit {
		b.flush()
	}
}

func (b *Batcher[T]) flush() int {
	if len(b.batch) == 0 {
		return 0
	}
	batch := b.batch
	b.batch = make([]T, 0, b.nlimit)
	b.handler(batch)
	return len(batch)
}

func (b *Batcher[T]) drain() {
	for {
		select {
		case item, ok := <-b.channel:
			if !ok {
				b.flush()
				return
			}
			b.append(item)
		default:
			close(b.channel)
		}
	}
}

func (b *Batcher[T]) Stop() {
	close(b.stop)
	<-b.done
}
