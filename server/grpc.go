package server

import (
	pb "github.com/504dev/logr-go-client/protos/gen/go"
	"github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/config"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"net"
)

type server struct {
	pb.UnimplementedLogRpcServer
}

func (s *server) Push(ctx context.Context, lrp *pb.LogRpcPackage) (*pb.Response, error) {
	var lp types.LogPackage
	lp.FromProto(lrp)
	HandleLog(&LogPackageMeta{
		LogPackage: &lp,
		Protocol:   "grpc",
		Size:       proto.Size(lrp),
	})
	return &pb.Response{}, nil
}

func MustListenGRPC() {
	if err := ListenGRPC(); err != nil {
		panic(err)
	}
}
func ListenGRPC() error {
	addr := config.Get().Bind.Grpc
	if addr == "" {
		return nil
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterLogRpcServer(s, &server{})
	return s.Serve(listener)
}
