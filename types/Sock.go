package types

import (
	_types "github.com/504dev/logr-go-client/types"
	"golang.org/x/net/websocket"
	"sync"
)

type Sock struct {
	sync.RWMutex
	SockId          string         `json:"sock_id"`
	Listeners       map[string]int `json:"listeners"`
	Paused          bool           `json:"paused"`
	*User           `json:"user"`
	*Filter         `json:"filter"`
	*Claims         `json:"claims"`
	*websocket.Conn `json:"conn"`
}

func (s *Sock) SendLog(lg *_types.Log) error {
	m := SockMessage{
		Path:    "/log",
		Payload: lg,
	}
	return websocket.JSON.Send(s.Conn, m)
}

func (s *Sock) AddListener(path string) {
	s.Lock()
	if s.Listeners == nil {
		s.Listeners = make(map[string]int)
	}
	s.Listeners[path] += 1
	s.Unlock()
}

func (s *Sock) RemoveListener(path string) {
	s.Lock()
	s.Listeners[path] -= 1
	s.Unlock()
}

func (s *Sock) SetFilter(f *Filter) {
	s.Lock()
	s.Filter = f
	s.Unlock()
}

func (s *Sock) SetPaused(state bool) {
	s.Lock()
	s.Paused = state
	s.Unlock()
}

type SockMessage struct {
	Action  string      `json:"action,omitempty"`
	Path    string      `json:"path"`
	Payload interface{} `json:"payload"`
}
