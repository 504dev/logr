package controllers

import (
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/user"
	"github.com/504dev/logr/models/ws"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type WsController struct{}

func (wc WsController) Index(c *gin.Context) {
	handler := websocket.Handler(wc.Reader)
	handler.ServeHTTP(c.Writer, c.Request)
}

func (wc WsController) Reader(w *websocket.Conn) {
	cfg := w.Config()
	query := cfg.Location.Query()
	tokenstring := query.Get("token")
	sockId := query.Get("sock_id")
	paused := query.Get("paused") == "true"

	if tokenstring == "" || sockId == "" {
		return
	}

	claims := &types.Claims{}
	tkn, err := claims.ParseWithClaims(tokenstring, config.Get().GetJwtSecret())

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
		SockId: sockId,
		User:   usr,
		Conn:   w,
		Paused: paused,
		Claims: claims,
	}
	ws.SockMap.Add(sock)

	for {
		var m types.SockMessage

		if err := websocket.JSON.Receive(w, &m); err != nil {
			Logger.Error("websocket.JSON.Receive: userId=%v, sockId=%v, err=%v", usr.Id, sockId, err)
			ws.SockMap.Delete(usr.Id, sockId)
			break
		}

		switch m.Action {
		case "subscribe":
			sock.AddListener(m.Path)
		case "unsubscribe":
			sock.RemoveListener(m.Path)
		case "pause":
			paused := m.Payload.(bool)
			sock.SetPaused(paused)
		}

		Logger.Debug("Received: %v %v", sockId, m)
	}
}
