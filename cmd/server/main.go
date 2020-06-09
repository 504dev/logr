package main

import (
	"github.com/504dev/logr/clickhouse"
	"github.com/504dev/logr/config"
	"github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/count"
	"github.com/504dev/logr/models/log"
	"github.com/504dev/logr/models/ws"
	"github.com/504dev/logr/mysql"
	"github.com/504dev/logr/server"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"time"
)

func main() {
	color.NoColor = false
	config.Init()
	clickhouse.Init()
	mysql.Init()
	logger.Init()
	log.RunQueue()
	count.RunQueue()
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
			logger.Logger.Info(ws.SockMap.Info())
		}
	})()
	HandleExit()
}

func HandleExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sig := <-c
	logger.Logger.Warn("Exit with code: %v", sig)
	log.StopQueue()
	count.StopQueue()
	os.Exit(0)
}
