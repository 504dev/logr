package logger

import (
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/models/dashkey"
	lgc "github.com/504dev/logr-go-client"
	"strconv"
)

type logger struct {
	*lgc.Logger
	*lgc.Counter
	Gin *lgc.Writter
}

func (lg *logger) Init() {
	conf := lgc.Config{
		Udp: config.Get().Bind.Udp,
	}
	dk, err := dashkey.GetById(1)
	if dk != nil {
		conf.DashId = dk.DashId
		conf.PublicKey = dk.PublicKey
		conf.PrivateKey = dk.PrivateKey
	}
	lg.Logger, _ = conf.NewLogger("main.log")
	lg.Counter, _ = conf.NewCounter("main.cnt")
	if err != nil {
		lg.Error(err)
	}
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

var Logger = logger{}
