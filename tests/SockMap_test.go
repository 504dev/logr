package tests

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/types"
	"github.com/504dev/logr/types/jwtservice"
	"github.com/504dev/logr/types/sockmap"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type writter struct {
	journal chan string
}

func (w writter) Drain() []string {
	var result []string
	for {
		select {
		case v := <-w.journal:
			result = append(result, v)
		default:
			return result
		}
	}
}

func (w writter) Write(p []byte) (int, error) {
	w.journal <- string(p)
	return 0, nil
}

func (w writter) Close() error {
	return nil
}

func TestSockMap(t *testing.T) {
	t.Parallel()
	sm := sockmap.NewSockMap()
	journal := make(chan string, 100)
	w := &writter{journal: journal}
	claims := &jwtservice.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}

	user1 := &types.User{Id: 1}
	user2 := &types.User{Id: 3}

	sock1 := &sockmap.Sock{User: user1, SockId: "sock_1", Conn: w, Claims: claims}
	sock2 := &sockmap.Sock{User: user1, SockId: "sock_2", Conn: w, Claims: claims}
	sock3 := &sockmap.Sock{User: user2, SockId: "sock_3", Conn: w, Claims: claims}

	sm.Register(sock1)
	sm.Register(sock2)
	sm.Register(sock3)

	sock1.AddListener("/log")
	sock2.AddListener("/log")
	sock3.AddListener("/log")

	sm.SetFilter(sock1.User.Id, sock1.SockId, &types.Filter{DashId: 1})
	sm.SetFilter(sock2.User.Id, sock2.SockId, &types.Filter{DashId: 2})
	sm.SetFilter(sock3.User.Id, sock3.SockId, &types.Filter{DashId: 1, Level: "error"})

	var result []string

	type testCase struct {
		*_types.Log
		len    int
		result []string
	}
	tests := []testCase{
		{
			&_types.Log{DashId: 1, Level: "info", Message: "hello"},
			1,
			nil,
		},
		{
			&_types.Log{DashId: 1, Level: "error", Message: "drop database"},
			2,
			nil,
		},
	}

	for _, tc := range tests {
		sm.Push(tc.Log)
		result = w.Drain()
		assert.Equal(t, tc.len, len(result), "Unexpected journal size")
		if tc.result != nil {
			assert.Equal(t, tc.result, result, "Unexpected journal content")
		}
	}

	close(w.journal)
}
