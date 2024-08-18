package http

import (
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/server/http/router"
	"github.com/504dev/logr/types/jwtservice"
	"github.com/504dev/logr/types/sockmap"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net"
	"net/http"
	"time"
)

type HTTPServer struct {
	jwtService *jwtservice.JwtService
	sockMap    *sockmap.SockMap
	engine     *gin.Engine
	server     *http.Server
	listener   net.Listener
}

func NewHTTPServer(
	addr string,
	sockMap *sockmap.SockMap,
	jwtService *jwtservice.JwtService,
	repos *repo.Repos,
) (*HTTPServer, error) {
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

	return &HTTPServer{
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

func (srv *HTTPServer) Engine() *gin.Engine {
	return srv.engine
}

func (srv *HTTPServer) Listen() error {
	return srv.server.Serve(srv.listener)
}

func (srv *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return srv.server.Shutdown(ctx)
}
