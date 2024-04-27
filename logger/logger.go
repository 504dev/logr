package logger

import (
	logr "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/config"
	"github.com/504dev/logr/models/dashkey"
	"github.com/504dev/logr/types"
	"strconv"
)

func createConfig(dashId int) (*logr.Config, error) {
	conf := logr.Config{
		Udp:      config.Get().Bind.Udp,
		Grpc:     config.Get().Bind.Grpc,
		NoCipher: false,
	}
	dk, err := dashkey.GetById(dashId)
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
	Support, _ = conf.NewLogger("support.log")
	conf.DefaultSystemCounter()
	conf.DefaultProcessCounter()
	gin, _ := conf.NewLogger("gin.log")
	GinWriter = gin.CustomWriter(func(log *logr.Log) {
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
	if config.Get().DemoDash != nil {
		go Demo()
	}
}

var Support *logr.Logger
var Logger *logr.Logger
var GinWriter *logr.Writer
