package cachify

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

const defaultExpiration = 15 * time.Second
const cleanupInterval = 10 * time.Minute

var store = cache.New(defaultExpiration, cleanupInterval)

const WARMING_TIME = 3 * time.Second
const PREWARM_TIME = 3 * time.Second
const PREWARM_ENABLED = true

func Cachify(key string, f func() (any, error), expired time.Duration) (any, error) {
	warmingKey := getWarmingKey(key)

	entry, exp, exist := store.GetWithExpiration(key)
	if exist {
		if PREWARM_ENABLED {
			if diff := -time.Since(exp); diff < PREWARM_TIME {
				if _, exist := store.Get(warmingKey); !exist {
					go func() {
						_, _ = fetchWithWarming(key, f, expired+diff) // PREWARM
					}()
				}
			}
		}

		return entry, nil
	}

	if _, locked := store.Get(warmingKey); locked {
		time.Sleep(time.Second)
		return Cachify(key, f, expired)
	}

	return fetchWithWarming(key, f, expired)
}

func fetchWithWarming(key string, f func() (any, error), expired time.Duration) (any, error) {
	warmingKey := getWarmingKey(key)

	store.Set(warmingKey, true, WARMING_TIME)
	at := time.Now()
	dash, err := f()
	if time.Since(at) < WARMING_TIME {
		store.Delete(warmingKey)
	}

	if err != nil {
		return nil, err
	}

	store.Set(key, dash, expired)

	return dash, nil
}

func getWarmingKey(key string) string {
	return fmt.Sprintf("%v:warming", key)
}
