package sockmap

import (
	"encoding/json"
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/types"
	"golang.org/x/net/websocket"
	"io"
	"sync"
	"time"
)

type SockSession struct {
	Paused        bool           `json:"paused"`
	Listeners     map[string]int `json:"listeners"`
	*types.Filter `json:"filter"`
}

type Sock struct {
	mu            sync.RWMutex
	store         SessionStore
	SockId        string         `json:"sock_id"`
	Session       *SockSession   `json:"session"`
	Conn          io.WriteCloser `json:"conn"`
	JwtToken      string         `json:"jwt_token"`
	*types.User   `json:"user"`
	*types.Claims `json:"claims"`
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
	go s.store.Set(s.SockId, s.Session)
}

func (s *Sock) LoadSession() {
	s.Session = &SockSession{}
	if sess, err := s.store.Get(s.SockId); err == nil && sess != nil {
		s.Session = sess
	}
}

func (s *Sock) SetStore(store SessionStore) {
	s.store = store
}

func (s *Sock) SendLog(lg *_types.Log) error {
	msg := SockMessage{
		Path:    "/log",
		Payload: lg,
	}
	switch conn := s.Conn.(type) {
	case *websocket.Conn:
		return websocket.JSON.Send(conn, msg)
	default:
		bytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		_, err = s.Conn.Write(bytes)
		return err
	}
}

func (s *Sock) IsExpired() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.RegisteredClaims.ExpiresAt.Before(time.Now())
}

func (s *Sock) HasListener(path string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Session.Listeners != nil && s.Session.Listeners[path] != 0
}

func (s *Sock) AddListener(path string) {
	s.mu.Lock()
	if s.Session.Listeners == nil {
		s.Session.Listeners = make(map[string]int)
	}
	s.Session.Listeners[path] += 1
	s.mu.Unlock()
}

func (s *Sock) RemoveListener(path string) {
	s.mu.Lock()
	s.Session.Listeners[path] -= 1
	s.mu.Unlock()
}

func (s *Sock) GetFilter() *types.Filter {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Session.Filter
}
func (s *Sock) SetFilter(f *types.Filter) {
	s.mu.Lock()
	s.Session.Filter = f
	s.mu.Unlock()
	go s.store.Set(s.SockId, s.Session)
}

func (s *Sock) IsPaused() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Session.Paused
}
func (s *Sock) SetPaused(state bool) {
	s.mu.Lock()
	s.Session.Paused = state
	s.mu.Unlock()
}

func (s *Sock) Delete() error {
	go s.store.Del(s.SockId)
	return s.Conn.Close()
}
