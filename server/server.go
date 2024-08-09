package server

import (
	_types "github.com/504dev/logr-go-client/types"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"golang.org/x/net/context"
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
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewLogServer(
	httpaddr string,
	udpAddr string,
	grpcAddr string,
	logStorage LogStorage,
	countStorage CountStorage,
) (result *LogServer, err error) {
	ch := make(chan *LogPackageMeta)
	ctx, cancel := context.WithCancel(context.Background())
	result = &LogServer{
		channel:      ch,
		joiner:       types.NewLogPackageJoiner(time.Second, 5),
		logStorage:   logStorage,
		countStorage: countStorage,
		ctx:          ctx,
		cancel:       cancel,
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

func (srv *LogServer) processLogs() {
	defer close(srv.channel)
	for {
		select {
		case <-srv.ctx.Done():
			return
		case meta := <-srv.channel:
			srv.handleLog(meta)
		}
	}

}

func (srv *LogServer) Run() {
	go srv.processLogs()
	go srv.udpServer.Listen()
	go func() {
		if err := srv.httpServer.Listen(); err != nil {
			Logger.Error(err)
		}
	}()
	go func() {
		if err := srv.grpcServer.Listen(); err != nil {
			panic(err)
		}
	}()
}

func (srv *LogServer) Stop() {
	srv.cancel()
	srv.udpServer.Stop()
	_ = srv.grpcServer.Stop()
	_ = srv.httpServer.Stop()
}
