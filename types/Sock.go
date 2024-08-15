package types

import (
	_types "github.com/504dev/logr-go-client/types"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

type session struct {
	Paused    bool           `json:"paused"`
	Listeners map[string]int `json:"listeners"`
	*Filter   `json:"filter"`
}

type Sock struct {
	sync.RWMutex
	session         *session
	store           SessionStore
	SockId          string `json:"sock_id"`
	JwtToken        string `json:"jwt_token"`
	*User           `json:"user"`
	*Claims         `json:"claims"`
	*websocket.Conn `json:"conn"` // TODO interface
}

type SockMessage struct {
	Action  string      `json:"action,omitempty"`
	Path    string      `json:"path"`
	Payload interface{} `json:"payload"`
}

func (s *Sock) HandleMessage(msg *SockMessage) {
	switch msg.Action {
	case "subscribe":
		s.AddListener(msg.Path)
	case "unsubscribe":
		s.RemoveListener(msg.Path)
	case "pause":
		paused := msg.Payload.(bool)
		s.SetPaused(paused)
	}
	go s.store.Set(s.SockId, s.session)
}

func (s *Sock) SendLog(lg *_types.Log) error {
	m := SockMessage{
		Path:    "/log",
		Payload: lg,
	}
	return websocket.JSON.Send(s.Conn, m)
}

func (s *Sock) IsExpired() bool {
	s.RLock()
	defer s.RUnlock()
	return s.RegisteredClaims.ExpiresAt.Before(time.Now())
}

func (s *Sock) HasListener(path string) bool {
	s.RLock()
	defer s.RUnlock()
	return s.session.Listeners != nil && s.session.Listeners[path] != 0
}

func (s *Sock) AddListener(path string) {
	s.Lock()
	if s.session.Listeners == nil {
		s.session.Listeners = make(map[string]int)
	}
	s.session.Listeners[path] += 1
	s.Unlock()
}

func (s *Sock) RemoveListener(path string) {
	s.Lock()
	s.session.Listeners[path] -= 1
	s.Unlock()
}

func (s *Sock) GetFilter() *Filter {
	s.RLock()
	defer s.RUnlock()
	return s.session.Filter
}
func (s *Sock) SetFilter(f *Filter) {
	s.Lock()
	s.session.Filter = f
	s.Unlock()
	go s.store.Set(s.SockId, s.session)
}

func (s *Sock) IsPaused() bool {
	s.RLock()
	defer s.RUnlock()
	return s.session.Paused
}
func (s *Sock) SetPaused(state bool) {
	s.Lock()
	s.session.Paused = state
	s.Unlock()
}
