package types

import (
	"golang.org/x/net/websocket"
	"sync"
)

type Sock struct {
	sync.RWMutex
	SockId          string `json:"sock_id"`
	*User           `json:"user"`
	*Filter         `json:"filter"`
	*websocket.Conn `json:"conn"`
}

func (s *Sock) SendLog(lg *Log) error {
	m := SockMessage{
		Path:    "/log",
		Payload: lg,
	}
	return websocket.JSON.Send(s.Conn, m)
}

func (s *Sock) SetFilter(f *Filter) {
	s.Lock()
	s.Filter = f
	s.Unlock()
}

type SockMap map[int]map[string]*Sock

func (sm SockMap) PushLog(lg *Log) int {
	cnt := 0
	for _, m := range sm {
		for _, s := range m {
			if s.Filter != nil && s.Filter.Match(lg) {
				err := s.SendLog(lg)
				if err != nil {
					sm.Delete(s.User.Id, s.SockId)
				}
				cnt += 1
			}
		}
	}
	return cnt
}

func (sm SockMap) Get(userId int, sockId string) *Sock {
	if _, ok := sm[userId]; ok {
		return sm[userId][sockId]
	}
	return nil
}

func (sm SockMap) SetFilter(userId int, sockId string, filter *Filter) bool {
	s := sm.Get(userId, sockId)
	if s != nil {
		s.SetFilter(filter)
		return true
	}
	return false
}

func (sm SockMap) Set(s *Sock) {
	if _, ok := sm[s.User.Id]; !ok {
		sm[s.User.Id] = make(map[string]*Sock)
	}
	sm[s.User.Id][s.SockId] = s
}

func (sm SockMap) Delete(userId int, uid string) bool {
	if _, ok := sm[userId]; !ok {
		return false
	}
	if _, ok := sm[userId][uid]; !ok {
		return false
	}
	delete(sm[userId], uid)
	return true
}

type SockMessage struct {
	Path    string      `json:"path"`
	Payload interface{} `json:"payload"`
}
