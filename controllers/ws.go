package controllers

import (
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/user"
	"github.com/504dev/kidlog/models/ws"
	"github.com/504dev/kidlog/types"
	"github.com/dgrijalva/jwt-go"
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
	token := query.Get("token")
	sockId := query.Get("sock_id")

	if token == "" || sockId == "" {
		return
	}

	claims := &types.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().OAuth.JwtSecret), nil
	})

	logger.Debug(claims)
	logger.Debug(tkn, err)

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
	}
	ws.SockMap.Set(sock)

	for {
		var m types.SockMessage

		if err := websocket.JSON.Receive(w, &m); err != nil {
			logger.Error("websocket.JSON.Receive: %v", err)
			ws.SockMap.Delete(usr.Id, sockId)
			break
		}
		switch m.Action {
		case "subscribe":
			sock.AddListener(m.Path)
		case "unsubscribe":
			sock.RemoveListener(m.Path)
		}

		logger.Debug("Received: %v", m)

		//if err := websocket.JSON.Send(w, m); err != nil {
		//	logger.Error("websocket.JSON.Send: %v", err)
		//	ws.SockMap.Delete(usr.Id, sockId)
		//	break
		//}
	}
}
