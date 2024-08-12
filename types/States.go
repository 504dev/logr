package types

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
	"math/rand"
	"time"
)

type States struct {
	data cmap.ConcurrentMap[string, string]
}

func NewStates() *States {
	return &States{
		data: cmap.New[string](),
	}
}

func (s *States) Pop(state string) (string, bool) {
	v, ok := s.data.Get(state)
	s.data.Remove(state)
	return v, ok
}

func (s *States) Push(v string) string {
	state := fmt.Sprintf("%v_%v", time.Now().Nanosecond(), rand.Int())
	s.data.Set(state, v)
	return state
}

func (s *States) Set(k string, v string) {
	s.data.Set(k, v)
}
