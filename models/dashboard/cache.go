package dashboard

import (
	"fmt"
	"github.com/504dev/logr/cachify"
	"github.com/504dev/logr/types"
)

func GetByIdCached(id int) (*types.Dashboard, error) {
	key := fmt.Sprintf("dashboard:id:%v", id)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return findOneByField("id", id)
	})
	return res.(*types.Dashboard), err
}
