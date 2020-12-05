package log

import (
	"fmt"
	"github.com/504dev/logr/cachify"
	"github.com/504dev/logr/types"
	"time"
)

func GetDashLognamesCached(dashId int) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("logs:lognames:%v", dashId)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return GetDashLognames(dashId)
	}, time.Minute)
	return res.([]*types.DashStatRow), err
}

func GetDashStatsCached(dashId int, logname string) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("logs:stats:%v:%v", dashId, logname)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return GetDashStats(dashId, logname)
	}, time.Minute)
	return res.([]*types.DashStatRow), err
}
