package ws

import (
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/types/jwtservice"
	"github.com/504dev/logr/types/sockmap"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"time"
)

type WsServer struct {
	repos      *repo.Repos
	jwtService *jwtservice.JwtService
	sockMap    *sockmap.SockMap
}

func NewWsServer(sockMap *sockmap.SockMap, jwtService *jwtservice.JwtService, repos *repo.Repos) *WsServer {
	return &WsServer{
		repos:      repos,
		jwtService: jwtService,
		sockMap:    sockMap,
	}
}

func (ws WsServer) Bind(e *gin.Engine) {
	e.GET("/ws", ws.Handshake)
}

func (ws WsServer) Handshake(ctx *gin.Context) {
	handler := websocket.Handler(ws.Stream)
	handler.ServeHTTP(ctx.Writer, ctx.Request)
}

func (ws WsServer) Stream(conn *websocket.Conn) {
	cfg := conn.Config()
	query := cfg.Location.Query()
	tokenstring := query.Get("token")
	sockId := query.Get("sock_id")

	if tokenstring == "" || sockId == "" {
		return
	}

	claims, tkn, err := ws.jwtService.ParseToken(tokenstring)

	Logger.Debug(claims)
	Logger.Debug(err, tkn)

	if err != nil || !tkn.Valid {
		return
	}

	usr, err := ws.repos.User.GetById(claims.Id)
	if err != nil {
		return
	}

	sock := &sockmap.Sock{
		SockId:   sockId,
		User:     usr,
		Conn:     conn,
		Claims:   claims,
		JwtToken: tokenstring,
	}
	ws.sockMap.Register(sock)

	for {
		var msg sockmap.SockMessage

		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			Logger.Error("websocket.JSON.Receive: userId=%v, sockId=%v, err=%v", usr.Id, sockId, err)
			ws.sockMap.Unregister(sock)

			break
		}

		sock.HandleMessage(&msg)

		Logger.Debug("Received: %v %v", sockId, msg)
	}
}

func (ws WsServer) Info() {
	const interval = 10 * time.Second

	go func() {
		for {
			time.Sleep(interval)
			Logger.Info("SockMap %v", ws.sockMap)
		}
	}()
}
