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
	*time.Ticker
	list    [][]interface{}
	flushed time.Time
}

func NewQueue(c *QueueConfig) *Queue {
	return &Queue{
		QueueConfig: c,
		flushed:     time.Now(),
		list:        make([][]interface{}, 0, c.FlushCount),
	}
}

func (q *Queue) Run() {
	q.Ticker = time.NewTicker(q.FlushInterval)
	go (func() {
		for {
			<-q.Ticker.C
			err := q.Flush()
			if err != nil {
				Logger.Error(err)
			}
		}
	})()
}

func (q *Queue) Stop() error {
	q.Ticker.Stop()
	return q.Flush()
}

func (q *Queue) Push(values []interface{}) {
	q.Lock()
	q.list = append(q.list, values)
	q.Unlock()
	if len(q.list) >= q.FlushCount {
		q.Flush()
	}
}

func (q *Queue) Flush() error {
	if len(q.list) == 0 {
		return nil
	}
	q.Lock()
	batch := q.list[0:]
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
