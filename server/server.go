package server

import (
	"github.com/504dev/logr-go-client/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"
)

type LogPackageMeta struct {
	*types.LogPackage
	Protocol string
	Size     int
}

type LogStorage interface {
	Store(entry *types.Log) error
}

type LogServer struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	udpServer  *UdpServer
	logChannel chan *LogPackageMeta
	storage    LogStorage
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewLogServer(udpaddr string) (*LogServer, error) {
	ch := make(chan *LogPackageMeta)
	udp, err := NewUdpServer(udpaddr, ch)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &LogServer{
		udpServer:  udp,
		logChannel: ch,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

func (srv *LogServer) Handle() {
	defer close(srv.logChannel)
	for {
		select {
		case <-srv.ctx.Done():
			return
		case msg := <-srv.logChannel:
			Handle(msg.LogPackage, msg.Protocol, msg.Size)
		}
	}

}
func (srv *LogServer) ListenUDP() {
	srv.udpServer.Listen()
}

func (srv *LogServer) Stop() {
	srv.cancel()
	srv.udpServer.Stop()
}
