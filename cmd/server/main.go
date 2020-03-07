package main

import (
	"encoding/json"
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/log"
	"github.com/504dev/kidlog/models/ws"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/server"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.Init()
	clickhouse.Init()
	mysql.Init()
	log.RunQueue()
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
	HandleExit()
}

func HandleExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sig := <-c
	logger.Warn("Exit with code: %v", sig)
	err := log.StopQueue()
	if err != nil {
		logger.Warn("Exit error: %v", err)
	}
	os.Exit(0)
}
