package log

import (
	"fmt"
	"github.com/504dev/logr/libs/cachify"
	"github.com/504dev/logr/types"
	"time"
)

func (repo *LogRepo) StatsByDashboardCached(dashId int) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("logs:lognames:%v", dashId)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return repo.StatsByDashboard(dashId)
	}, time.Minute)
	if err != nil {
		return nil, err
	}
	return res.([]*types.DashStatRow), nil
}

func (repo *LogRepo) StatsByLognameCached(dashId int, logname string) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("logs:stats:%v:%v", dashId, logname)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return repo.StatsByLogname(dashId, logname)
	}, time.Minute)
	if err != nil {
		return nil, err
	}

	return res.([]*types.DashStatRow), nil
}
