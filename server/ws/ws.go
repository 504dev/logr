package ws

import (
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/repo/user"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"time"
)

type WsServer struct {
	jwtService *types.JwtService
	sockmap    *types.SockMap
}

func NewWsServer(sockmap *types.SockMap, jwtService *types.JwtService) *WsServer {
	return &WsServer{
		jwtService: jwtService,
		sockmap:    sockmap,
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

	usr, err := user.GetById(claims.Id)
	if err != nil {
		return
	}

	sock := &types.Sock{
		SockId:   sockId,
		User:     usr,
		Conn:     conn,
		Claims:   claims,
		JwtToken: tokenstring,
	}
	ws.sockmap.Register(sock)

	for {
		var msg types.SockMessage

		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			Logger.Error("websocket.JSON.Receive: userId=%v, sockId=%v, err=%v", usr.Id, sockId, err)
			ws.sockmap.Unregister(sock)
			break
		}

		sock.HandleMessage(&msg)

		Logger.Debug("Received: %v %v", sockId, msg)
	}
}

func (ws WsServer) Info() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			Logger.Info("SockMap %v", ws.sockmap)
		}
	}()
}
