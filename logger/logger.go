package logger

import (
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/models/dashkey"
	"github.com/504dev/kidlog/types"
	lgc "github.com/504dev/logr-go-client"
	"strconv"
)

func createConfig(dashId int) (*lgc.Config, error) {
	conf := lgc.Config{
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
	*lgc.Logger
	*lgc.Counter
	Gin *lgc.Writter
}

func (lg *loggerT) Init() {
	conf := lgc.Config{
		Udp: config.Get().Bind.Udp,
	}
	dk, _ := dashkey.GetById(types.DashboardSystemId)
	if dk != nil {
		conf.DashId = dk.DashId
		conf.PublicKey = dk.PublicKey
		conf.PrivateKey = dk.PrivateKey
	}
	lg.Logger, _ = conf.NewLogger("main.log")
	lg.Counter, _ = conf.NewCounter("main.cnt")
	gin, _ := conf.NewLogger("gin.log")
	lg.Gin = gin.CustomWritter(func(log *lgc.Log) {
		codestr := log.Message[38:41]
		code, _ := strconv.Atoi(codestr)
		if code >= 400 && code <= 499 {
			log.Level = lgc.LevelWarn
		} else if code >= 500 && code <= 599 {
			log.Level = lgc.LevelError
		}
	})
}

var Logger = loggerT{}
