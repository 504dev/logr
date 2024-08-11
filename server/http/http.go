package http

import (
	"github.com/504dev/logr/server/http/router"
	"github.com/504dev/logr/types"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpServer struct {
	iam     *types.AuthService
	sockmap *types.SockMap
	engine  *gin.Engine
	server  *http.Server
}

func NewHttpServer(addr string, sockmap *types.SockMap, iam *types.AuthService) (*HttpServer, error) {
	frontend := func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	}

	engine := router.NewRouter(sockmap, iam)
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
		iam:     iam,
		sockmap: sockmap,
		engine:  engine,
		server: &http.Server{
			Addr:    addr,
			Handler: engine,
		},
	}, nil
}

func (srv *HttpServer) Engine() *gin.Engine {
	return srv.engine
}

func (srv *HttpServer) Listen() error {
	return srv.server.ListenAndServe()
}

func (srv *HttpServer) Stop() error {
	return srv.server.Close()
}
