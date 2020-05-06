package logger

import (
	"github.com/504dev/kidlog/types"
	"github.com/fatih/color"
	"time"
)

func (lg *loggerT) Demo() {
	conf, _ := createConfig(types.DashboardDemoId)
	go (func() {
		l, _ := conf.NewLogger("starwars.log")
		i := 0
		for {
			c := crowls[i%len(crowls)]
			l.Warn(color.New(color.Bold).SprintFunc()(c.title))
			time.Sleep(3 * time.Second)
			for _, t := range c.text {
				l.Info(t)
				time.Sleep(1 * time.Second)
			}
			i += 1
		}
	})()
	go crypto(conf)
}
