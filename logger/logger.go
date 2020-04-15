package logger

import (
	"github.com/504dev/kidlog/config"
	lgc "github.com/504dev/logr-go-client"
	"strconv"
)

type logger struct {
	*lgc.Logger
	*lgc.Counter
	Gin *lgc.Writter
}

func (lg *logger) Init() {
	var options = config.Get().Logger
	var conf = lgc.Config{
		Udp:        options.Udp,
		DashId:     options.DashId,
		PrivateKey: options.PrivateKey,
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

var Logger = logger{}
