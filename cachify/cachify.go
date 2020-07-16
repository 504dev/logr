package cachify

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

var store = cache.New(15*time.Second, 10*time.Minute)

func Cachify(key string, f func() (interface{}, error), expired time.Duration) (interface{}, error) {
	entry, exist := store.Get(key)
	if exist {
		return entry, nil
	}
	warming := fmt.Sprintf("%v:warming", key)
	entry, exist = store.Get(warming)
	if exist {
		time.Sleep(time.Second)
		return Cachify(key, f, expired)
	}

	store.Set(warming, true, 3*time.Second)
	dash, err := f()
	store.Delete(warming)

	if err != nil {
		return nil, err
	}

	store.Set(key, dash, expired)

	return dash, nil
}
