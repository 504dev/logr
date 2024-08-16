package udp

import (
	"encoding/json"
	_types "github.com/504dev/logr-go-client/types"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"golang.org/x/net/context"
	"net"
	"sync"
)

type UdpServer struct {
	conn   *net.UDPConn
	ch     chan<- *types.LogPackageMeta
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

func NewUdpServer(addr string, ch chan<- *types.LogPackageMeta) (*UdpServer, error) {
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
		conn:   udpconn,
		ch:     ch,
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}, nil
}

func (srv *UdpServer) Listen() error {
	if srv == nil {
		return nil
	}

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, 10) // limit concurrent connections

	buf := make([]byte, 65536)
	for {
		if srv.ctx.Err() != nil {
			break
		}

		size, _, err := srv.conn.ReadFromUDP(buf)
		if err != nil {
			Logger.Error("UDP read error: %v", err)
			continue
		}

		data := make([]byte, size)
		copy(data, buf[:size])

		semaphore <- struct{}{}
		wg.Add(1)

		go func() {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			lp := _types.LogPackage{}
			if err := json.Unmarshal(data, &lp); err != nil {
				Logger.Error("UDP parse json error: %v\n%v", err, string(data))
				return
			}

			srv.ch <- &types.LogPackageMeta{
				LogPackage: &lp,
				Protocol:   "udp",
				Size:       size,
			}
		}()
	}
	wg.Wait()
	close(srv.done)

	return nil
}

func (srv *UdpServer) Stop() error {
	if srv == nil {
		return nil
	}
	srv.cancel()
	<-srv.done
	return srv.conn.Close()
}
