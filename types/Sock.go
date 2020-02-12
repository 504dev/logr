package types

import "golang.org/x/net/websocket"

type Sock struct {
	Uid string
	*User
	*Filter
	*websocket.Conn
}

type SockMap map[int]Sock

type SockMessage struct {
	Path    string      `json:"path"`
	Payload interface{} `json:"payload"`
}
