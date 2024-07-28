package ws

import (
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"sync"
	"time"
)

var once sync.Once
var sockmap *types.SockCMap

func GetSockMap() *types.SockCMap {
	once.Do(func() {
		sockmap = types.NewSockCMap()
		if connstring := config.Get().Redis; connstring != "" {
			sss, err := types.NewRedisSessionStore(connstring, 0, time.Hour)
			if err != nil {
				Logger.Error("cannot create redis session store: %v", err)
				return
			}
			sockmap.SetSessionStore(sss)
		}
	})
	return sockmap
}
