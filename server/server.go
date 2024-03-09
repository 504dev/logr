package server

import (
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func MustListenHTTP() {
	if err := ListenHTTP(); err != nil {
		panic(err)
	}
}
func ListenHTTP() error {
	gin.ForceConsoleColor()

	gin.DefaultWriter = io.MultiWriter(os.Stdout, GinWriter)

	// TODO react
	r := NewRouter()
	r.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
	r.GET("/", frontend)
	r.GET("/demo", frontend)
	r.GET("/login", frontend)
	r.GET("/jwt/:token", frontend)
	r.GET("/dashboards", frontend)
	r.GET("/dashboard/*rest", frontend)
	r.GET("/policy", frontend)
	r.GET("/support", frontend)

	return r.Run(config.Get().Bind.Http)
}

func frontend(c *gin.Context) {
	c.File("./frontend/dist/index.html")
}
