package logger

import (
	"github.com/504dev/logr/types"
)

func Demo() {
	conf, _ := createConfig(types.DashboardDemoId)
	go author(conf)
	go starwars(conf)
	go crypto(conf)
}
