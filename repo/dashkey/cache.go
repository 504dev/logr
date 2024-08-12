package dashkey

import (
	"fmt"
	"github.com/504dev/logr/libs/cachify"
	"github.com/504dev/logr/types"
	"time"
)

func (repo *DashboardKeyRepo) GetByPubCached(pub string) (*types.DashKey, error) {
	key := fmt.Sprintf("dashkey:pub:%v", pub)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return repo.GetByPub(pub)
	}, 15*time.Second)
	return res.(*types.DashKey), err
}
