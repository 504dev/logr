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

	clickhouse.Init(config.GetCommandLineArgs().Retries)
	mysql.Init(config.GetCommandLineArgs().Retries)

	repos := repo.GetRepos()

	logger.Init(repos.DashboardKey)

	logServer, err := server.NewLogServer(
		config.Get().Bind.Http,
		config.Get().Bind.Udp,
		config.Get().Bind.Grpc,
		config.Get().Redis,
		config.Get().GetJwtSecret,
		repos,
	)
	if err != nil {
		panic(err)
	}
	logServer.Run()

	// Exit & Graceful Shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-exit
	ts := time.Now()
	logger.Logger.Warn("Exit with code: %v", sig)
	logServer.Stop()
	logger.Logger.Debug(time.Since(ts))
}
