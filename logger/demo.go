package logger

import (
	"github.com/504dev/logr/logger/ai"
	"github.com/504dev/logr/types"
)

func Demo() {
	conf, _ := createConfig(types.DashboardDemoId)
	go ai.Run(conf)
	go starwars(conf)
	go crypto(conf)
}
