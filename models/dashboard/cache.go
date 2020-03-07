package dashboard

import (
	"fmt"
	"github.com/504dev/kidlog/types"
	cache "github.com/patrickmn/go-cache"
	"time"
)

var c = cache.New(15*time.Second, 10*time.Minute)

func cachify(key string, f func() (interface{}, error)) (interface{}, error) {
	entry, exist := c.Get(key)
	if exist {
		return entry, nil
	}
	warming := fmt.Sprintf("%v:warming", key)
	entry, exist = c.Get(warming)
	if exist {
		time.Sleep(time.Second)
		return cachify(key, f)
	}

	c.Set(warming, true, 3*time.Second)
	dash, err := f()
	c.Delete(warming)

	if err != nil {
		return nil, err
	}

	c.SetDefault(key, dash)

	return dash, nil
}

func GetByIdCached(id int) (*types.Dashboard, error) {
	key := fmt.Sprintf("dashboard:id:%v", id)
	res, err := cachify(key, func() (interface{}, error) {
		return findOneByField("id", id)
	})
	return res.(*types.Dashboard), err
}

func GetByPubCached(pub string) (*types.Dashboard, error) {
	key := fmt.Sprintf("dashboard:pub:%v", pub)
	res, err := cachify(key, func() (interface{}, error) {
		return findOneByField("public_key", pub)
	})
	return res.(*types.Dashboard), err
}
