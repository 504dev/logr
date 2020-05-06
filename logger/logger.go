package logger

import (
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/models/dashkey"
	"github.com/504dev/kidlog/types"
	logr "github.com/504dev/logr-go-client"
	"strconv"
)

func createConfig(dashId int) (*logr.Config, error) {
	conf := logr.Config{
		Udp: config.Get().Bind.Udp,
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

type loggerT struct {
	*logr.Logger
	*logr.Counter
	Gin *logr.Writter
}

func (lg *loggerT) Init() {
	conf, _ := createConfig(types.DashboardSystemId)
	lg.Logger, _ = conf.NewLogger("main.log")
	lg.Counter, _ = conf.NewCounter("main.cnt")
	gin, _ := conf.NewLogger("gin.log")
	lg.Gin = gin.CustomWritter(func(log *logr.Log) {
		codestr := log.Message[38:41]
		code, _ := strconv.Atoi(codestr)
		if code >= 400 && code <= 499 {
			log.Level = logr.LevelWarn
		} else if code >= 500 && code <= 599 {
			log.Level = logr.LevelError
		}
	})
	go lg.Demo()
}

var Logger = &loggerT{}
