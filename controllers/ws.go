package controllers

import (
	"fmt"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"io"
)

type WsController struct{}

func (wc WsController) Index(c *gin.Context) {
	handler := websocket.Handler(wc.Reader2)
	handler.ServeHTTP(c.Writer, c.Request)
}

type message struct {
	// the json tag means this will serialize as a lowercased field
	Message string `json:"message"`
}

func (_ WsController) Reader2(ws *websocket.Conn) {
	for {
		// allocate our container struct
		var m message

		// receive a message using the codec
		if err := websocket.JSON.Receive(ws, &m); err != nil {
			logger.Error(err)
			break
		}

		logger.Debug("Received message:", m.Message)

		// send a response
		m2 := message{"Thanks for the message!"}
		if err := websocket.JSON.Send(ws, m2); err != nil {
			logger.Error(err)
			break
		}
	}
}

func (_ WsController) Reader(ws *websocket.Conn) {
	cfg := ws.Config()
	query := cfg.Location.Query()
	token := query.Get("token")
	fmt.Println(cfg)
	claims := &types.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().OAuth.JwtSecret), nil
	})
	fmt.Println(tkn, err)
	fmt.Println(claims.Username, claims.Role, claims.Id)
	io.Copy(ws, ws)
}
