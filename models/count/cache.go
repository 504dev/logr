package count

import (
	"fmt"
	"github.com/504dev/logr/cachify"
	"github.com/504dev/logr/types"
	"time"
)

func StatsByDashboardCached(dashId int) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("counts:lognames:%v", dashId)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return StatsByDashboard(dashId)
	}, time.Minute)
	if err != nil {
		return nil, err
	}
	return res.([]*types.DashStatRow), nil
}

func StatsByLognameCached(dashId int, logname string) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("counts:stats:%v:%v", dashId, logname)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return StatsByLogname(dashId, logname)
	}, time.Minute)
	if err != nil {
		return nil, err
	}
	return res.([]*types.DashStatRow), nil
}
