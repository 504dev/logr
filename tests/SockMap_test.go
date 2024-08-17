package tests

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type writter struct {
	journal chan string
}

func (w writter) Write(p []byte) (n int, err error) {
	w.journal <- string(p)
	return 0, nil
}
func (w writter) Close() error {
	return nil
}

func TestSockMap(t *testing.T) {
	sm := types.NewSockMap()
	journal := make(chan string, 100)
	w := &writter{journal: journal}
	claims := &types.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}

	user1 := &types.User{Id: 1}
	user2 := &types.User{Id: 3}

	sock1 := &types.Sock{User: user1, SockId: "sock_1", Conn: w, Claims: claims}
	sock2 := &types.Sock{User: user1, SockId: "sock_2", Conn: w, Claims: claims}
	sock3 := &types.Sock{User: user2, SockId: "sock_3", Conn: w, Claims: claims}

	sm.Register(sock1)
	sm.Register(sock2)
	sm.Register(sock3)

	sock1.AddListener("/log")
	sock2.AddListener("/log")
	sock3.AddListener("/log")

	sm.SetFilter(sock1.User.Id, sock1.SockId, &types.Filter{DashId: 1})
	sm.SetFilter(sock2.User.Id, sock2.SockId, &types.Filter{DashId: 2})
	sm.SetFilter(sock3.User.Id, sock3.SockId, &types.Filter{DashId: 1, Logname: "hello.log"})

	sm.Push(&_types.Log{DashId: 1, Logname: "hello.log"})

	time.Sleep(time.Second)
	close(w.journal)
	assert.Equal(t, len(w.journal), 2, "Unexpected journal size")
	for v := range w.journal {
		t.Log(v)
	}
}
