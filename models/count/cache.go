package count

import (
	"fmt"
	"github.com/504dev/logr/cachify"
	"github.com/504dev/logr/types"
	"time"
)

func GetDashLognamesCached(dashId int) ([]*types.DashStatRow, error) {
	key := fmt.Sprintf("counts:lognames:%v", dashId)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return GetDashLognames(dashId)
	}, time.Minute)
	return res.([]*types.DashStatRow), err
}
