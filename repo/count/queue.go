package count

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/types"
)

func (repo *CountRepo) RunQueue() *CountRepo {
	repo.queue.Run()
	return repo
}

func (repo *CountRepo) StopQueue() error {
	return repo.queue.Stop()
}

func (repo *CountRepo) Store(c *_types.Count) error {
	cv := (*types.CountVector)(c)
	repo.queue.Push(cv.AsVector())
	return nil
}
