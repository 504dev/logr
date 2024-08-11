package dashboard

import (
	"fmt"
	"github.com/504dev/logr/libs/cachify"
	"github.com/504dev/logr/types"
	"time"
)

func GetByIdCached(id int) (*types.Dashboard, error) {
	key := fmt.Sprintf("dashboard:id:%v", id)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return findOneByField("id", id)
	}, 15*time.Second)
	return res.(*types.Dashboard), err
}
