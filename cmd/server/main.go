package main

import (
	"github.com/504dev/logr/clickhouse"
	"github.com/504dev/logr/config"
	"github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/count"
	"github.com/504dev/logr/models/log"
	"github.com/504dev/logr/mysql"
	"github.com/504dev/logr/server"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	color.NoColor = false
	args := config.Init()
	clickhouse.Init(args.Retries)
	mysql.Init(args.Retries)
	logger.Init()

	logStorage := log.NewLogStorage().RunQueue()
	countStorage := count.NewCountStorage().RunQueue()
	logServer, err := server.NewLogServer(
		config.Get().Bind.Http,
		config.Get().Bind.Udp,
		config.Get().Bind.Grpc,
		logStorage,
		countStorage,
	)
	if err != nil {
		panic(err)
	}
	logServer.Run()

	// Shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigchan
	logger.Logger.Warn("Exit with code: %v", sig)
	//logServer.Stop()
	_ = logStorage.StopQueue()
	_ = countStorage.StopQueue()
	os.Exit(0)
}
