package count

import (
	_types "github.com/504dev/logr-go-client/types"
)

func (repo *CountRepo) RunBatcher() *CountRepo {
	go repo.batcher.Run()
	return repo
}

func (repo *CountRepo) Store(count *_types.Count) error {
	repo.batcher.Push(count)
	return nil
}

func (repo *CountRepo) StopBatcher() error {
	repo.batcher.Stop()
	return nil
}
