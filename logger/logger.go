package logger

import (
	logr "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/config"
	"github.com/504dev/logr/logger/demo"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
)

func createConfig(dashId int) (*logr.Config, error) {
	conf := logr.Config{
		Udp:      config.Get().Bind.Udp,
		Grpc:     config.Get().Bind.Grpc,
		NoCipher: false,
	}
	dk, err := repo.GetRepos().DashboardKey.GetById(dashId)
	if err != nil {
		return nil, err
	}
	if dk != nil {
		conf.DashId = dk.DashId
		conf.PublicKey = dk.PublicKey
		conf.PrivateKey = dk.PrivateKey
	}
	return &conf, err
}

func Init() {
	conf, _ := createConfig(types.DashboardSystemId)
	Logger, _ = conf.NewLogger("main.log")
	_, _ = conf.DefaultSystemCounter()
	_, _ = conf.DefaultProcessCounter()

	gin.ForceConsoleColor()
	gin.DefaultWriter = io.MultiWriter(os.Stdout, GinWriter())

	if config.Get().DemoDash.Enabled {
		conf, _ := createConfig(types.DashboardDemoId)
		go demo.Run(conf, Logger)
	}
}

func GinWriter() *logr.Writer {
	ginlog, _ := Logger.Config.NewLogger("gin.log")
	return ginlog.CustomWriter(func(log *logr.Log) {
		codestr := log.Message[38:41]
		code, _ := strconv.Atoi(codestr)
		if code >= 400 && code <= 499 {
			log.Level = logr.Levels.Warn.String()
		} else if code >= 500 && code <= 599 {
			log.Level = logr.Levels.Error.String()
		}
		if code > 0 {
			log.Message = log.Message[28:]
		}
	})
}

var Logger *logr.Logger
