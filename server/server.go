package server

import (
	_types "github.com/504dev/logr-go-client/types"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"golang.org/x/sync/errgroup"
	"time"
)

type LogPackageMeta struct {
	*_types.LogPackage
	Protocol string
	Size     int
}

type LogStorage interface {
	Store(*_types.Log) error
}

type CountStorage interface {
	Store(*_types.Count) error
}

type LogServer struct {
	httpServer   *HttpServer
	grpcServer   *GrpcServer
	udpServer    *UdpServer
	channel      chan *LogPackageMeta
	joiner       *types.LogPackageJoiner
	logStorage   LogStorage
	countStorage CountStorage
	done         chan struct{}
}

func NewLogServer(
	httpaddr string,
	udpAddr string,
	grpcAddr string,
	logStorage LogStorage,
	countStorage CountStorage,
) (result *LogServer, err error) {
	ch := make(chan *LogPackageMeta)
	result = &LogServer{
		channel:      ch,
		joiner:       types.NewLogPackageJoiner(time.Second, 5),
		logStorage:   logStorage,
		countStorage: countStorage,
		done:         make(chan struct{}),
	}
	if udpAddr != "" {
		result.udpServer, err = NewUdpServer(udpAddr, ch)
		if err != nil {
			return nil, err
		}
	}
	if grpcAddr != "" {
		result.grpcServer, err = NewGrpcServer(grpcAddr, ch)
		if err != nil {
			return nil, err
		}
	}

	result.httpServer, err = NewHttpServer(httpaddr)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (srv *LogServer) processChannel() {
	for meta := range srv.channel {
		srv.handleLog(meta)
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
