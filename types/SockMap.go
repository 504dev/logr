package types

import (
	"encoding/json"
	_types "github.com/504dev/logr-go-client/types"
	"sync"
	"time"
)

type SockMap struct {
	sync.RWMutex
	data map[int]map[string]*Sock
}

func (sm *SockMap) PushLog(lg *_types.Log) int {
	cnt := 0
	sm.Lock()
	now := time.Now().Unix()
	for _, m := range sm.data {
		for _, s := range m {
			if s.Filter == nil || s.Paused || s.Listeners == nil || s.Listeners["/log"] == 0 {
				continue
			}
			if s.ExpiresAt < now {
				sm.delete(s.User.Id, s.SockId)
				continue
			}
			if s.Filter.Match(lg) {
				err := s.SendLog(lg)
				if err != nil {
					sm.delete(s.User.Id, s.SockId)
				}
				cnt += 1
			}
		}
	}
	sm.Unlock()
	return cnt
}

func (sm *SockMap) SetFilter(userId int, sockId string, filter *Filter) bool {
	s := sm.Get(userId, sockId)
	if s != nil {
		s.SetFilter(filter)
		return true
	}
	return false
}

func (sm *SockMap) SetPaused(userId int, sockId string, state bool) bool {
	s := sm.Get(userId, sockId)
	if s != nil {
		s.SetPaused(state)
		return true
	}
	return false
}

func (sm *SockMap) Get(userId int, sockId string) *Sock {
	sm.RLock()
	defer sm.RUnlock()
	if _, ok := sm.data[userId]; ok {
		return sm.data[userId][sockId]
	}
	return nil
}

func (sm *SockMap) init() {
	if sm.data == nil {
		sm.data = make(map[int]map[string]*Sock)
	}
}

func (sm *SockMap) Set(s *Sock) {
	sm.Lock()
	sm.init()
	if _, ok := sm.data[s.User.Id]; !ok {
		sm.data[s.User.Id] = make(map[string]*Sock)
	}
	sm.data[s.User.Id][s.SockId] = s
	sm.Unlock()
}

func (sm *SockMap) delete(userId int, uid string) bool {
	if _, ok := sm.data[userId]; !ok {
		return false
	}
	if _, ok := sm.data[userId][uid]; !ok {
		return false
	}
	sm.data[userId][uid].Close()
	delete(sm.data[userId], uid)
	return true
}

func (sm *SockMap) Delete(userId int, uid string) bool {
	sm.Lock()
	flag := sm.delete(userId, uid)
	sm.Unlock()
	return flag
}

func (sm *SockMap) Info() string {
	j, _ := json.Marshal(sm.data)
	return string(j)
}
