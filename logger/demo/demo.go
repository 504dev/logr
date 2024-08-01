package demo

import (
	logr "github.com/504dev/logr-go-client"
	"github.com/504dev/logr/logger/demo/ai"
)

func Run(conf *logr.Config, mainlog *logr.Logger) {
	go ai.Run(conf)
	go starwars(conf)
	go crypto(conf, mainlog)
}
