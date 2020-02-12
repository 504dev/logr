package controllers

import (
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/logger"
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

type message struct {
	Message string `json:"message"`
}

func (wc WsController) Reader(ws *websocket.Conn) {
	claims, tkn, err := wc.EnsureJWT(ws)
	logger.Debug(claims)
	logger.Debug(tkn, err)

	if err != nil || !tkn.Valid {
		return
	}

	for {
		var m message

		if err := websocket.JSON.Receive(ws, &m); err != nil {
			logger.Error(err)
			break
		}

		logger.Debug("Received message: %v", m.Message)

		m2 := message{"Thanks for the message!"}
		if err := websocket.JSON.Send(ws, m2); err != nil {
			logger.Error(err)
			break
		}
	}
}

func (_ WsController) EnsureJWT(ws *websocket.Conn) (*types.Claims, *jwt.Token, error) {
	cfg := ws.Config()
	logger.Debug(cfg)
	query := cfg.Location.Query()
	token := query.Get("token")
	claims := &types.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().OAuth.JwtSecret), nil
	})
	return claims, tkn, err
}
