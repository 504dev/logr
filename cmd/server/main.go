package main

import (
	"encoding/json"
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/ws"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/server"
	"time"
)

func main() {
	config.Init()
	clickhouse.Init()
	mysql.Init()
	go (func() {
		err := server.ListenUDP()
		if err != nil {
			panic(err)
		}
	})()
	go (func() {
		err := server.ListenHTTP()
		if err != nil {
			panic(err)
		}
	})()
	go (func() {
		for {
			time.Sleep(10 * time.Second)
			j, _ := json.Marshal(ws.SockMap)
			logger.Info(string(j))
		}
	})()
	select {}
}
