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

type SockMap struct {
	clients *cmap.ConcurrentMap[uid, *cmap.ConcurrentMap[string, *Sock]]
	SockSessionStore
}

func NewSockMap() *SockMap {
	clients := cmap.NewStringer[uid, *cmap.ConcurrentMap[string, *Sock]]()
	sm := SockMap{
		clients:          &clients,
		SockSessionStore: &MemorySessionStore{},
	}
	return &sm
}

func (sm *SockMap) SetSessionStore(store SockSessionStore) {
	sm.SockSessionStore = store
}

func (sm *SockMap) Register(s *Sock) {
	s.SockSession = &SockSession{}
	s.SockSessionStore = sm.SockSessionStore
	// load session
	session, err := sm.SockSessionStore.Get(s.SockId)
	if err == nil && session != nil {
		s.SockSession = session
	}
	go sm.SockSessionStore.Set(s.SockId, session)
	sm.add(s)
}

func (sm *SockMap) Unregister(s *Sock) {
	go sm.SockSessionStore.Del(s.SockId)
	sm.delete(s)
}

func (sm *SockMap) Push(lg *_types.Log) int {
	cnt := 0
	now := time.Now()
	// TODO index socks by lg.DashId
	for user := range sm.clients.IterBuffered() {
		for sock := range user.Val.IterBuffered() {
			s := sock.Val
			sFilter := s.GetFilter()
			if sFilter == nil || s.IsPaused() || !s.HasListener("/log") {
				continue
			}
			if s.Claims.ExpiresAt.Before(now) {
				sm.Unregister(s)
				continue
			}
			if sFilter.Match(lg) {
				if err := s.SendLog(lg); err != nil {
					sm.Unregister(s)
				}
				cnt += 1
			}
		}
	}
	return cnt
}

func (sm *SockMap) SetFilter(userId int, sockId string, filter *Filter) bool {
	if s := sm.GetSock(userId, sockId); s != nil {
		s.SetFilter(filter)
		return true
	}
	return false
}

func (sm *SockMap) GetSocks(userId int) *cmap.ConcurrentMap[string, *Sock] {
	us, _ := sm.clients.Get(uid(userId))
	return us
}

func (sm *SockMap) GetSock(userId int, sockId string) *Sock {
	if us := sm.GetSocks(userId); us != nil {
		sock, _ := us.Get(sockId)
		return sock
	}
	return nil
}

func (sm *SockMap) add(s *Sock) {
	us := sm.GetSocks(s.User.Id)
	if us != nil {
		us.Set(s.SockId, s)
		return
	}
	tmp := cmap.New[*Sock]()
	tmp.Set(s.SockId, s)
	sm.clients.Set(uid(s.User.Id), &tmp)
}

func (sm *SockMap) delete(s *Sock) bool {
	if us, ok := sm.clients.Get(uid(s.User.Id)); ok {
		if s, ok := us.Get(s.SockId); ok {
			_ = s.Close()
			us.Remove(s.SockId)
			return true
		}
	}
	return false
}

func (sm *SockMap) String() string {
	j, _ := json.Marshal(sm.clients)
	return string(j)
}
