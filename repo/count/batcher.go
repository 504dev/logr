package count

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/dbs/clickhouse/batcher"
	. "github.com/504dev/logr/logger"
	"time"
)

func (repo *CountRepo) StartBatcher() {
	const batchSize = 1000
	const batchInterval = time.Second
	repo.batcher = batcher.NewBatcher(batchSize, batchInterval, func(batch []*_types.Count) {
		err := repo.BatchInsert(batch)
		Logger.InfoErr(err, "Batch insert %v %v", len(batch), err)
	})

	go repo.batcher.Run()
}

func (repo *CountRepo) Store(count *_types.Count) error {
	if repo.batcher != nil {
		repo.batcher.Push(count)
		return nil
	}
	return repo.Insert(count)
}

func (repo *CountRepo) StopBatcher() {
	if repo.batcher != nil {
		repo.batcher.Stop()
	}
}
