package logger

import (
	"errors"
	logr "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/config"
	"github.com/504dev/logr/logger/demo"
	"github.com/504dev/logr/repo/interfaces"
	"github.com/504dev/logr/types"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
)

var ErrNoDashKey = errors.New("no dashkey provided")

func createConfig(repo interfaces.DashboardKeyRepo, id int) *logr.Config {
	dashkey, err := repo.GetById(id)
	if err != nil {
		panic(err)
	}
	if dashkey == nil {
		panic(ErrNoDashKey)
	}
	return &logr.Config{
		Udp:        config.Get().Bind.UDP,
		Grpc:       config.Get().Bind.GRPC,
		DashId:     dashkey.DashId,
		PublicKey:  dashkey.PublicKey,
		PrivateKey: dashkey.PrivateKey,
		NoCipher:   false,
	}
}

func Init(repo interfaces.DashboardKeyRepo) {
	color.NoColor = false

	conf := createConfig(repo, types.DASHKEY_SYSTEM_ID)
	Logger, _ = conf.NewLogger("main.log")
	_, _ = conf.DefaultSystemCounter()
	_, _ = conf.DefaultProcessCounter()

	gin.ForceConsoleColor()
	gin.DefaultWriter = io.MultiWriter(os.Stdout, GinWriter())

	if config.Get().DemoDash.Enabled {
		conf := createConfig(repo, types.DASHKEY_DEMO_ID)
		go demo.Start(conf, Logger)
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
