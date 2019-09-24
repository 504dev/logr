package main

import (
	"../../clickhouse"
	"../../config"
	"../../types"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	clickhouse.Init()

	r := gin.Default()
	r.GET("/filter", func(c *gin.Context) {
		c.JSON(200, types.Logs{})
	})

	r.Run(config.Get().Bind)
}
