package server

import (
	_types "github.com/504dev/logr-go-client/types"
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
	Store(entry *_types.Log) error
}

type LogServer struct {
	httpServer *HttpServer
	grpcServer *GrpcServer
	udpServer  *UdpServer
	logChannel chan *LogPackageMeta
	joiner     *types.LogPackageJoiner
	storage    LogStorage
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewLogServer(httpaddr string, udpAddr string, grpcAddr string) (result *LogServer, err error) {
	ch := make(chan *LogPackageMeta)
	ctx, cancel := context.WithCancel(context.Background())
	result = &LogServer{
		logChannel: ch,
		joiner:     types.NewLogPackageJoiner(time.Second, 5),
		ctx:        ctx,
		cancel:     cancel,
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
	defer close(srv.logChannel)
	for {
		select {
		case <-srv.ctx.Done():
			return
		case meta := <-srv.logChannel:
			srv.handleLog(meta)
		}
	}

}

func (srv *LogServer) Run() {
	go srv.processLogs()
	go srv.udpServer.Listen()
	go func() {
		if err := srv.httpServer.Listen(); err != nil {
			panic(err)
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
	srv.httpServer.Stop()
	srv.udpServer.Stop()
	_ = srv.grpcServer.Stop()
}
