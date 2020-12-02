package types

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type States struct {
	Data map[string]string
	sync.Mutex
}

func (s States) Get(state string) (string, bool) {
	s.Lock()
	v, ok := s.Data[state]
	delete(s.Data, state)
	s.Unlock()
	return v, ok
}
func (s States) Insert(v string) string {
	state := fmt.Sprintf("%v_%v", time.Now().Nanosecond(), rand.Int())
	s.Lock()
	s.Data[state] = v
	s.Unlock()
	return state
}

func (s States) Set(k string, v string) {
	s.Lock()
	s.Data[k] = v
	s.Unlock()
}
