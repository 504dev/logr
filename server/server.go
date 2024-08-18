package server

import (
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/server/grpc"
	"github.com/504dev/logr/server/http"
	"github.com/504dev/logr/server/udp"
	"github.com/504dev/logr/server/ws"
	"github.com/504dev/logr/types"
	"github.com/504dev/logr/types/jwtservice"
	"github.com/504dev/logr/types/sockmap"
	"golang.org/x/sync/errgroup"
	nethttp "net/http"
	"time"
)

type LogServer struct {
	httpServer *http.HttpServer
	wsServer   *ws.WsServer
	grpcServer *grpc.GrpcServer
	udpServer  *udp.UdpServer
	jwtService *jwtservice.JwtService
	sockMap    *sockmap.SockMap
	channel    chan *types.LogPackageMeta
	joiner     *types.LogPackageJoiner
	repos      *repo.Repos
	done       chan struct{}
}

func NewLogServer(
	httpAddr string,
	udpAddr string,
	grpcAddr string,
	redisAddr string,
	jwtSecretFunc func() string,
	repos *repo.Repos,
) (*LogServer, error) {
	var err error
	var udpServer *udp.UdpServer
	var grpcServer *grpc.GrpcServer
	var httpServer *http.HttpServer
	var wsServer *ws.WsServer

	channel := make(chan *types.LogPackageMeta)
	if udpAddr != "" {
		udpServer, err = udp.NewUdpServer(udpAddr, channel)
		if err != nil {
			return nil, err
		}
	}
	if grpcAddr != "" {
		grpcServer, err = grpc.NewGrpcServer(grpcAddr, channel)
		if err != nil {
			return nil, err
		}
	}

	jwtService := jwtservice.NewJwtService(jwtSecretFunc)

	sockMap := sockmap.NewSockMap()
	if redisAddr != "" {
		store, err := sockmap.NewRedisSessionStore(redisAddr, time.Hour, time.Second)
		if err != nil {
			return nil, err
		}
		sockMap.SetSessionStore(store)
	}

	httpServer, err = http.NewHttpServer(httpAddr, sockMap, jwtService, repos)
	if err != nil {
		return nil, err
	}

	wsServer = ws.NewWsServer(sockMap, jwtService, repos)
	wsServer.Bind(httpServer.Engine())
	wsServer.Info()

	return &LogServer{
		udpServer:  udpServer,
		grpcServer: grpcServer,
		httpServer: httpServer,
		wsServer:   wsServer,
		jwtService: jwtService,
		sockMap:    sockMap,
		channel:    channel,
		joiner:     types.NewLogPackageJoiner(time.Second, 5),
		repos:      repos,
		done:       make(chan struct{}),
	}, nil
}

func (srv *LogServer) handleLoop() {
	for meta := range srv.channel {
		go srv.handle(meta)
	}
}

func (srv *LogServer) Run() {
	go func() {
		if err := srv.httpServer.Listen(); err != nil {
			if err == nethttp.ErrServerClosed {
				Logger.Warn(err)
				return
			}
			panic(err)
		}
	}()

	// start recieving log packages
	var wg errgroup.Group
	wg.Go(srv.udpServer.Listen)
	wg.Go(srv.grpcServer.Listen)
	go func() {
		if err := wg.Wait(); err != nil {
			Logger.Warn(err)
		}
		close(srv.channel) // writing to srv.channel is complete
	}()

	// start handling
	srv.handleLoop()
	close(srv.done) // reading from srv.channel completed
}

func (srv *LogServer) Stop() error {
	var wg errgroup.Group
	wg.Go(srv.udpServer.Stop)
	wg.Go(srv.grpcServer.Stop)
	wg.Go(srv.httpServer.Stop)
	if err := wg.Wait(); err != nil {
		return err
	}
	srv.repos.Stop()
	<-srv.done
	return nil
}
