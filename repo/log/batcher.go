package log

import (
	_types "github.com/504dev/logr-go-client/types"
)

func (repo *LogRepo) RunBatcher() *LogRepo {
	go repo.batcher.Run()
	return repo
}

func (repo *LogRepo) Store(log *_types.Log) error {
	repo.batcher.Push(log)
	return nil
}

func (repo *LogRepo) StopBatcher() error {
	repo.batcher.Stop()
	return nil
}
