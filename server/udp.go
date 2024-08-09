package server

import (
	"encoding/json"
	_types "github.com/504dev/logr-go-client/types"
	. "github.com/504dev/logr/logger"
	"golang.org/x/net/context"
	"net"
)

type UdpServer struct {
	conn   *net.UDPConn
	ch     chan *LogPackageMeta
	ctx    context.Context
	cancel context.CancelFunc
}

func NewUdpServer(addr string, ch chan *LogPackageMeta) (*UdpServer, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	udpconn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &UdpServer{
		udpconn,
		ch,
		ctx,
		cancel,
	}, nil
}

func (srv *UdpServer) Listen() {
	defer srv.conn.Close()
	buf := make([]byte, 65536)
	for {
		if srv.ctx.Err() != nil {
			return
		}
		size, _, err := srv.conn.ReadFromUDP(buf)
		if err != nil {
			Logger.Error("UDP read error: %v", err)
			continue
		}

		data := make([]byte, size)
		copy(data, buf[:size])

		go func() {
			lp := _types.LogPackage{}
			if err := json.Unmarshal(data, &lp); err != nil {
				Logger.Error("UDP parse json error: %v\n%v", err, string(data))
				return
			}

			srv.ch <- &LogPackageMeta{
				LogPackage: &lp,
				Protocol:   "udp",
				Size:       size,
			}
		}()
	}
}

func (srv *UdpServer) Stop() {
	srv.cancel()
}
