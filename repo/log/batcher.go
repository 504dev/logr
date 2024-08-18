package log

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/dbs/clickhouse/batcher"
	. "github.com/504dev/logr/logger"
	"time"
)

func (repo *LogRepo) StartBatcher() {
	const batchSize = 1000
	const batchInterval = time.Second / 2
	repo.batcher = batcher.NewBatcher(batchSize, batchInterval, func(batch []*_types.Log) {
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
