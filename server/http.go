package server

import (
	"context"
	. "github.com/504dev/logr/logger"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
)

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(addr string) (*HttpServer, error) {
	gin.ForceConsoleColor()

	gin.DefaultWriter = io.MultiWriter(os.Stdout, GinWriter())

	frontend := func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	}

	engine := NewRouter()
	engine.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
	engine.GET("/", frontend)
	engine.GET("/demo", frontend)
	engine.GET("/login", frontend)
	engine.GET("/jwt/:token", frontend)
	engine.GET("/dashboards", frontend)
	engine.GET("/dashboard/*rest", frontend)
	engine.GET("/policy", frontend)
	engine.GET("/support", frontend)

	return &HttpServer{
		server: &http.Server{
			Addr:    addr,
			Handler: engine,
		},
	}, nil
}

func (srv *HttpServer) Listen() error {
	return srv.server.ListenAndServe()
}

func (srv *HttpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return srv.server.Shutdown(ctx)
}
