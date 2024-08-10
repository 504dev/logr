package server

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(addr string) (*HttpServer, error) {
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
	return srv.server.Close()
}
