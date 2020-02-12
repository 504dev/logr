package controllers

import (
	"fmt"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"io"
)

type WsController struct{}

func (wc WsController) Index(c *gin.Context) {
	handler := websocket.Handler(wc.Reader)
	handler.ServeHTTP(c.Writer, c.Request)
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
