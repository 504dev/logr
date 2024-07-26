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
	args := config.Init()
	clickhouse.Init(args.Retries)
	mysql.Init(args.Retries)
	logger.Init()
	log.RunQueue()   // TODO graceful shutdown
	count.RunQueue() // TODO graceful shutdown
	go server.MustListenUDP()
	go server.MustListenGRPC()
	go server.MustListenHTTP()
	go server.MustListenPROM()
	go (func() {
		for {
			time.Sleep(10 * time.Second)
			logger.Logger.Info("ws.SockMap %v", ws.SockMap.String())
		}
	})()
	HandleExit()
}

func HandleExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sig := <-c
	logger.Logger.Warn("Exit with code: %v", sig)
	_ = log.StopQueue()
	_ = count.StopQueue()
	os.Exit(0)
}
