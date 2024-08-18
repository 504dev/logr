package grpc

import (
	pb "github.com/504dev/logr-go-client/protos/gen/go"
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"net"
)

type logRpcService struct {
	pb.UnimplementedLogRpcServer
	ch chan<- *types.LogPackageMeta
}

func (s *logRpcService) Push(_ context.Context, lrp *pb.LogRpcPackage) (*pb.Response, error) {
	var lp _types.LogPackage

	lp.FromProto(lrp)

	s.ch <- &types.LogPackageMeta{
		LogPackage: &lp,
		Protocol:   "grpc",
		Size:       proto.Size(lrp),
	}

	return &pb.Response{}, nil
}

type GRPCServer struct {
	grpcServer *grpc.Server
	listener   net.Listener
	service    *logRpcService
}

func NewGRPCServer(addr string, ch chan<- *types.LogPackageMeta) (*GRPCServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &GRPCServer{
		grpcServer: grpc.NewServer(),
		listener:   listener,
		service:    &logRpcService{ch: ch},
	}, nil
}

func (s *GRPCServer) Listen() error {
	if s == nil {
		return nil
	}
	pb.RegisterLogRpcServer(s.grpcServer, s.service)
	return s.grpcServer.Serve(s.listener)
}

func (s *GRPCServer) Stop() error {
	if s == nil {
		return nil
	}
	s.grpcServer.GracefulStop()
	return nil
}
