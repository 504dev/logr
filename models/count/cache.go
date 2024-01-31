package count

import (
	"fmt"
	"github.com/504dev/logr/cachify"
	"github.com/504dev/logr/types"
	"time"
)

func GetStatsDashboardCached(dashId int) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("counts:lognames:%v", dashId)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return GetStatsDashboard(dashId)
	}, time.Minute)
	return res.([]*types.DashStatRow), err
}

func GetStatsLognameCached(dashId int, logname string) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("counts:stats:%v:%v", dashId, logname)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return GetStatsLogname(dashId, logname)
	}, time.Minute)
	return res.([]*types.DashStatRow), err
}
