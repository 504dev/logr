package dashkey

import (
	"fmt"
	"github.com/504dev/logr/libs/cachify"
	"github.com/504dev/logr/types"
	"time"
)

func (repo *DashboardKeyRepo) GetByPubCached(pub string) (*types.DashKey, error) {
	const cacheTime = 15 * time.Second
	cacheKey := fmt.Sprintf("dashkey:pub:%v", pub)
	res, err := cachify.Cachify(cacheKey, func() (interface{}, error) {
		return repo.GetByPub(pub)
	}, cacheTime)
	return res.(*types.DashKey), err
}
