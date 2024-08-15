package logger

import (
	"errors"
	logr "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/config"
	"github.com/504dev/logr/logger/demo"
	"github.com/504dev/logr/types"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
)

func createConfig(dashkey *types.DashKey) *logr.Config {
	if dashkey == nil {
		panic(errors.New("no dashkey provided"))
	}
	return &logr.Config{
		Udp:        config.Get().Bind.Udp,
		Grpc:       config.Get().Bind.Grpc,
		DashId:     dashkey.DashId,
		PublicKey:  dashkey.PublicKey,
		PrivateKey: dashkey.PrivateKey,
		NoCipher:   false,
	}
}

func Init(dashkeys types.DashKeys) {
	color.NoColor = false

	conf := createConfig(dashkeys.Get(types.DASHKEY_SYSTEM_ID))
	Logger, _ = conf.NewLogger("main.log")
	_, _ = conf.DefaultSystemCounter()
	_, _ = conf.DefaultProcessCounter()

	gin.ForceConsoleColor()
	gin.DefaultWriter = io.MultiWriter(os.Stdout, GinWriter())

	if config.Get().DemoDash.Enabled {
		conf := createConfig(dashkeys.Get(types.DASHKEY_DEMO_ID))
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
