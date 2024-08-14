package count

import (
	_types "github.com/504dev/logr-go-client/types"
)

func (repo *CountRepo) RunQueue() *CountRepo {
	go repo.batcher.Run()
	return repo
}

func (repo *CountRepo) StopQueue() error {
	repo.batcher.Stop()
	return nil
}

func (repo *CountRepo) Store(count *_types.Count) error {
	repo.batcher.Push(count)
	return nil
}
