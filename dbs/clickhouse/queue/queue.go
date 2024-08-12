package queue

import (
	. "github.com/504dev/logr/logger"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type QueueConfig struct {
	*sqlx.DB
	FlushInterval time.Duration
	FlushCount    int
	Sql           string
}

type Queue struct {
	sync.Mutex
	*QueueConfig
	list    [][]interface{}
	flushed time.Time
	stop    chan struct{}
}

func NewQueue(c *QueueConfig) *Queue {
	return &Queue{
		QueueConfig: c,
		flushed:     time.Now(),
		list:        make([][]interface{}, 0, c.FlushCount),
		stop:        make(chan struct{}),
	}
}

func (q *Queue) Run() {
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

func (q *Queue) Stop() error {
	close(q.stop)
	return q.Flush()
}

func (q *Queue) Push(values []interface{}) {
	q.Lock()
	q.list = append(q.list, values)
	count := len(q.list)
	q.Unlock()
	if count >= q.FlushCount {
		q.Flush()
	}
}

func (q *Queue) Flush() error {
	q.Lock()
	if len(q.list) == 0 {
		q.Unlock()
		return nil
	}
	batch := q.list
	q.list = make([][]interface{}, 0, q.FlushCount)
	q.Unlock()

	tx, err := q.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(q.Sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range batch {
		_, err = stmt.Exec(v...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
