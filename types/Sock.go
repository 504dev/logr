package types

import (
	"golang.org/x/net/websocket"
	"sync"
)

type Sock struct {
	sync.RWMutex
	Uid             string `json:"uid"`
	*User           `json:"user"`
	*Filter         `json:"filter"`
	*websocket.Conn `json:"conn"`
}

func (s *Sock) SetFilter(f *Filter) {
	s.Lock()
	s.Filter = f
	s.Unlock()
}

type SockMap map[int]map[string]Sock

func (sm SockMap) Set(s *Sock) {
	if _, ok := sm[s.User.Id]; !ok {
		sm[s.User.Id] = make(map[string]Sock)
	}
	sm[s.User.Id][s.Uid] = *s
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
