package log

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/dbs/clickhouse/batcher"
	. "github.com/504dev/logr/logger"
	"time"
)

func (repo *LogRepo) StartBatcher() {
	repo.batcher = batcher.NewBatcher(1000, time.Second/2, func(batch []*_types.Log) {
		err := repo.BatchInsert(batch)
		Logger.InfoErr(err, "Batch insert %v %v", len(batch), err)
	})
	go repo.batcher.Run()
}

func (repo *LogRepo) Store(log *_types.Log) error {
	if repo.batcher != nil {
		repo.batcher.Push(log)
		return nil
	}
	return repo.Insert(log)
}

func (repo *LogRepo) StopBatcher() {
	if repo.batcher != nil {
		repo.batcher.Stop()
	}
}
