package types

import (
	"encoding/json"
	"fmt"
	_types "github.com/504dev/logr-go-client/types"
	cmap "github.com/orcaman/concurrent-map/v2"
	"time"
)

type UID int

func (id UID) String() string {
	return fmt.Sprintf("%v", int(id))
}

type SockCMap struct {
	data *cmap.ConcurrentMap[UID, *cmap.ConcurrentMap[string, *Sock]]
}

func (sm SockCMap) Init() *SockCMap {
	data := cmap.NewStringer[UID, *cmap.ConcurrentMap[string, *Sock]]()
	sm.data = &data
	return &sm
}

func (sm *SockCMap) PushLog(lg *_types.Log) int {
	cnt := 0
	now := time.Now().Unix()
	for user := range sm.data.IterBuffered() {
		for sock := range user.Val.IterBuffered() {
			s := sock.Val
			if s.Filter == nil || s.Paused || s.Listeners == nil || s.Listeners["/log"] == 0 {
				continue
			}
			if s.ExpiresAt < now {
				sm.Delete(s.User.Id, s.SockId)
				continue
			}
			if s.Filter.Match(lg) {
				if err := s.SendLog(lg); err != nil {
					sm.Delete(s.User.Id, s.SockId)
				}
				cnt += 1
			}
		}
	}
	return cnt
}

func (sm *SockCMap) SetFilter(userId int, sockId string, filter *Filter) bool {
	if s := sm.GetSock(userId, sockId); s != nil {
		s.SetFilter(filter)
		return true
	}
	return false
}

func (sm *SockCMap) SetPaused(userId int, sockId string, state bool) bool {
	if s := sm.GetSock(userId, sockId); s != nil {
		s.SetPaused(state)
		return true
	}
	return false
}

func (sm *SockCMap) GetSocks(userId int) *cmap.ConcurrentMap[string, *Sock] {
	us, _ := sm.data.Get(UID(userId))
	return us
}

func (sm *SockCMap) GetSock(userId int, sockId string) *Sock {
	if us := sm.GetSocks(userId); us != nil {
		sock, _ := us.Get(sockId)
		return sock
	}
	return nil
}

func (sm *SockCMap) Add(s *Sock) {
	us := sm.GetSocks(s.User.Id)
	if us == nil {
		tmp := cmap.New[*Sock]()
		tmp.Set(s.SockId, s)
		sm.data.Set(UID(s.User.Id), &tmp)
		return
	}
	us.Set(s.SockId, s)
}

func (sm *SockCMap) Delete(userId int, sockId string) bool {
	if us, ok := sm.data.Get(UID(userId)); ok {
		if s, ok := us.Get(sockId); ok {
			s.Close()
			us.Remove(sockId)
			return true
		}
	}
	return false
}

func (sm *SockCMap) String() string {
	j, _ := json.Marshal(sm.data)
	return string(j)
}
