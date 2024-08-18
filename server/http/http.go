package http

import (
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/server/http/router"
	"github.com/504dev/logr/types"
	"github.com/504dev/logr/types/sockmap"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net"
	"net/http"
	"time"
)

type HttpServer struct {
	jwtService *types.JwtService
	sockMap    *sockmap.SockMap
	engine     *gin.Engine
	server     *http.Server
	listener   net.Listener
}

func NewHttpServer(
	addr string,
	sockMap *sockmap.SockMap,
	jwtService *types.JwtService,
	repos *repo.Repos,
) (*HttpServer, error) {
	frontend := func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	}

	engine := router.NewRouter(sockMap, jwtService, repos)
	engine.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
	engine.GET("/", frontend)
	engine.GET("/demo", frontend)
	engine.GET("/login", frontend)
	engine.GET("/jwt/:token", frontend)
	engine.GET("/dashboards", frontend)
	engine.GET("/dashboard/*rest", frontend)
	engine.GET("/policy", frontend)
	engine.GET("/support", frontend)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &HttpServer{
		jwtService: jwtService,
		sockMap:    sockMap,
		engine:     engine,
		listener:   listener,
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
	return srv.server.Serve(srv.listener)
}

func (srv *HttpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return srv.server.Shutdown(ctx)
}
