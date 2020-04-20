package dashkey

import (
	"fmt"
	"github.com/504dev/kidlog/cachify"
	"github.com/504dev/kidlog/types"
)

func GetByPubCached(pub string) (*types.DashKey, error) {
	key := fmt.Sprintf("dashkey:pub:%v", pub)
	res, err := cachify.Cachify(key, func() (interface{}, error) {
		return findOneByField("public_key", pub)
	})
	return res.(*types.DashKey), err
}
