package server

import (
	_types "github.com/504dev/logr-go-client/types"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/server/grpc"
	"github.com/504dev/logr/server/http"
	"github.com/504dev/logr/server/udp"
	"github.com/504dev/logr/server/ws"
	"github.com/504dev/logr/types"
	"golang.org/x/sync/errgroup"
	"time"
)

type LogStorage interface {
	Store(*_types.Log) error
}

type CountStorage interface {
	Store(*_types.Count) error
}

type LogServer struct {
	httpServer   *http.HttpServer
	wsServer     *ws.WsServer
	grpcServer   *grpc.GrpcServer
	udpServer    *udp.UdpServer
	jwtService   *types.JwtService
	sockmap      *types.SockMap
	channel      chan *types.LogPackageMeta
	joiner       *types.LogPackageJoiner
	logStorage   LogStorage
	countStorage CountStorage
	repos        *repo.Repos
	done         chan struct{}
}

func NewLogServer(
	httpAddr string,
	udpAddr string,
	grpcAddr string,
	redisAddr string,
	jwtSecretFunc func() string,
	logStorage LogStorage,
	countStorage CountStorage,
	repos *repo.Repos,
) (*LogServer, error) {
	var err error
	var udpServer *udp.UdpServer
	var grpcServer *grpc.GrpcServer
	var httpServer *http.HttpServer
	var wsServer *ws.WsServer

	ch := make(chan *types.LogPackageMeta)
	if udpAddr != "" {
		udpServer, err = udp.NewUdpServer(udpAddr, ch)
		if err != nil {
			return nil, err
		}
	}
	if grpcAddr != "" {
		grpcServer, err = grpc.NewGrpcServer(grpcAddr, ch)
		if err != nil {
			return nil, err
		}
	}

	jwtService := types.NewJwtService(jwtSecretFunc)

	sockmap := types.NewSockMap()
	if redisAddr != "" {
		store, err := types.NewRedisSessionStore(redisAddr, time.Hour)
		if err != nil {
			return nil, err
		}
		sockmap.SetSessionStore(store)
	}

	httpServer, err = http.NewHttpServer(httpAddr, sockmap, jwtService, repos)
	if err != nil {
		return nil, err
	}

	wsServer = ws.NewWsServer(sockmap, jwtService, repos)
	wsServer.Bind(httpServer.Engine())
	wsServer.Info()

	return &LogServer{
		udpServer:    udpServer,
		grpcServer:   grpcServer,
		httpServer:   httpServer,
		wsServer:     wsServer,
		jwtService:   jwtService,
		sockmap:      sockmap,
		channel:      ch,
		joiner:       types.NewLogPackageJoiner(time.Second, 5),
		logStorage:   logStorage,
		countStorage: countStorage,
		repos:        repos,
		done:         make(chan struct{}),
	}, nil
}

func (srv *LogServer) processChannel() {
	for meta := range srv.channel {
		srv.handle(meta)
	}
}

func (srv *LogServer) Run() {
	go func() {
		srv.processChannel()
		close(srv.done) // reading from srv.channel completed
	}()
	go func() {
		if err := srv.httpServer.Listen(); err != nil {
			Logger.Warn(err)
		}
	}()
	var g errgroup.Group
	g.Go(func() error {
		srv.udpServer.Listen()
		return nil
	})
	g.Go(func() error {
		return srv.grpcServer.Listen()
	})
	go func() {
		if err := g.Wait(); err != nil {
			Logger.Warn(err)
		}
		close(srv.channel) // writing to srv.channel is complete
	}()
}

func (srv *LogServer) Stop() {
	srv.udpServer.Stop()
	_ = srv.grpcServer.Stop()
	_ = srv.httpServer.Stop()
	<-srv.done
}
