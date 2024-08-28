package main

import (
	"github.com/504dev/logr/config"
	"github.com/504dev/logr/dbs/clickhouse"
	"github.com/504dev/logr/dbs/mysql"
	"github.com/504dev/logr/logger"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.MustLoad()

	clickhouse.MustInit(config.GetCommandLineArgs().Retries)
	mysql.MustInit(config.GetCommandLineArgs().Retries)

	repos := repo.GetRepos()

	logger.Init(repos.DashboardKey)

	logServer, err := server.NewLogServer(
		config.Get().Bind.HTTP,
		config.Get().Bind.UDP,
		config.Get().Bind.GRPC,
		config.Get().Redis,
		config.Get().GetJwtSecret,
		repos,
	)
	if err != nil {
		panic(err)
	}

	go logServer.Run()

	// Exit & Graceful Shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-exit

	ts := time.Now()
	logger.Logger.Warn("Exit with code: %v", sig)
	err = logServer.Stop()
	if err != nil {
		logger.Logger.Error(err)
	}
	logger.Logger.Debug(time.Since(ts))
}
