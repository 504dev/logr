package server

import (
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func ListenHTTP() error {
	gin.ForceConsoleColor()

	gin.DefaultWriter = io.MultiWriter(os.Stdout, GinWritter)

	r := NewRouter()
	r.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
	r.Use(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	return r.Run(config.Get().Bind.Http)
}
