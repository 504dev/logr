package types

import (
	"encoding/json"
	_types "github.com/504dev/logr-go-client/types"
	cmap "github.com/orcaman/concurrent-map/v2"
	"strconv"
	"time"
)

type uid int

func (id uid) String() string {
	return strconv.Itoa(int(id))
}

func NewSockCMap() *SockCMap {
	data := cmap.NewStringer[uid, *cmap.ConcurrentMap[string, *Sock]]()
	return &SockCMap{&data}
}

type SockCMap struct {
	data *cmap.ConcurrentMap[uid, *cmap.ConcurrentMap[string, *Sock]]
}

func (sm *SockCMap) PushLog(lg *_types.Log) int {
	cnt := 0
	now := time.Now()
	// TODO index socks by lg.DashId
	for user := range sm.data.IterBuffered() {
		for sock := range user.Val.IterBuffered() {
			s := sock.Val
			sFilter := s.GetFilter()
			if sFilter == nil || s.IsPaused() || !s.HasListener("/log") {
				continue
			}
			if s.Claims.ExpiresAt.Before(now) {
				sm.Delete(s.User.Id, s.SockId)
				continue
			}
			if sFilter.Match(lg) {
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
	us, _ := sm.data.Get(uid(userId))
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
		sm.data.Set(uid(s.User.Id), &tmp)
		return
	}
	us.Set(s.SockId, s)
}

func (sm *SockCMap) Delete(userId int, sockId string) bool {
	if us, ok := sm.data.Get(uid(userId)); ok {
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
