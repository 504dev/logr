package logger

import (
	"github.com/504dev/logr/types"
	"github.com/fatih/color"
	"math/rand"
	"strings"
	"time"
)

func Demo() {
	conf, _ := createConfig(types.DashboardDemoId)
	go (func() {
		logger, _ := conf.NewLogger("starwars.log")
		for {
			c := crowls[rand.Intn(len(crowls))]
			logger.Warn(color.New(color.Bold).SprintFunc()(c.title))
			time.Sleep(time.Second)
			for _, t := range c.text {
				time.Sleep(333 * time.Millisecond)
				logger.Info(t)
				logger.Inc("count:letters", float64(len(t)))
				logger.Inc("count:words", float64(len(strings.Fields(t))))
				if strings.Contains(t, "Jedi") {
					logger.Inc("count:Jedi", 1)
				}
				if strings.Contains(t, "Leia") {
					logger.Inc("count:Leia", 1)
				}
			}
			time.Sleep(2 * time.Second)
		}
	})()
	go crypto(conf)
}
