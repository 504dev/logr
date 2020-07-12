package logger

import (
	"github.com/504dev/logr/types"
	"github.com/fatih/color"
	"strings"
	"time"
)

func Demo() {
	conf, _ := createConfig(types.DashboardDemoId)
	go (func() {
		logger, _ := conf.NewLogger("starwars.log")
		for i := 0; ; i += 1 {
			c := crowls[i%len(crowls)]
			logger.Warn(color.New(color.Bold).SprintFunc()(c.title))
			time.Sleep(3 * time.Second)
			for _, t := range c.text {
				logger.Info(t)
				logger.Inc("count:letters", float64(len(t)))
				logger.Inc("count:words", float64(len(strings.Fields(t))))
				if strings.Contains(t, "Jedi") {
					logger.Inc("count:Jedi", 1)
				}
				if strings.Contains(t, "Leia") {
					logger.Inc("count:Leia", 1)
				}

				time.Sleep(1 * time.Second)
			}

		}
	})()
	go crypto(conf)
}
