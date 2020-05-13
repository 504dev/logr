package main

import (
	"encoding/json"
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/config"
	. "github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/count"
	"github.com/504dev/kidlog/models/log"
	"github.com/504dev/kidlog/models/ws"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/server"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.Init()
	clickhouse.Init()
	mysql.Init()
	Logger.Init()
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
			Logger.Info(ws.SockMap.Info())
		}
	})()
	Logger.Debug(os.Environ())
	a := "I \033[31mlove\033[0m Stack Overflow"
	b := "I " + color.RedString("love") + " Stack Overflow"
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	Logger.Debug(a)
	Logger.Debug(b)
	Logger.Debug(string(ja))
	Logger.Debug(string(jb))
	Logger.Debug(string(ja) == string(jb))
	HandleExit()
}

func HandleExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sig := <-c
	Logger.Warn("Exit with code: %v", sig)
	log.StopQueue()
	count.StopQueue()
	os.Exit(0)
}
