package controllers

import (
	"encoding/json"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/user"
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

var SockMap = make(types.SockMap)

func (wc WsController) Reader(ws *websocket.Conn) {
	cfg := ws.Config()
	query := cfg.Location.Query()
	token := query.Get("token")
	uid := query.Get("uid")

	if token == "" || uid == "" {
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

	SockMap.Set(&types.Sock{
		Uid:  uid,
		User: usr,
		Conn: ws,
	})

	j, _ := json.MarshalIndent(SockMap, "", "    ")
	logger.Info(string(j))

	logger.Info(ws.IsClientConn())
	logger.Info(ws.IsServerConn())

	for {
		var m types.SockMessage

		if err := websocket.JSON.Receive(ws, &m); err != nil {
			logger.Error("websocket.JSON.Receive: %v", err)
			SockMap.Delete(usr.Id, uid)
			break
		}

		logger.Debug("Received payload: %v", m.Payload)

		if err := websocket.JSON.Send(ws, m); err != nil {
			logger.Error("websocket.JSON.Send: %v", err)
			SockMap.Delete(usr.Id, uid)
			break
		}
	}
}
