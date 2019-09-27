package main

import (
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/types"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	clickhouse.Init()

	r := gin.Default()
	r.GET("/filter", func(c *gin.Context) {
		c.JSON(200, types.Logs{})
	})

	r.Run(config.Get().Bind.Http)
}
