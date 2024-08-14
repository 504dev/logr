package count

import (
	"fmt"
	"github.com/504dev/logr/libs/cachify"
	"github.com/504dev/logr/types"
	"time"
)

func (repo *CountRepo) StatsByDashboardCached(dashId int) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("counts:lognames:%v", dashId)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return repo.StatsByDashboard(dashId)
	}, time.Minute)
	if err != nil {
		return nil, err
	}
	return res.([]*types.DashStatRow), nil
}

func (repo *CountRepo) StatsByLognameCached(dashId int, logname string) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("counts:stats:%v:%v", dashId, logname)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return repo.StatsByLogname(dashId, logname)
	}, time.Minute)
	if err != nil {
		return nil, err
	}
	return res.([]*types.DashStatRow), nil
}
