package udp

import (
	"encoding/json"
	_types "github.com/504dev/logr-go-client/types"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"net"
	"sync"
	"sync/atomic"
)

const bufferSize = 65536
const concurrentLimit = 10

type UDPServer struct {
	conn *net.UDPConn
	ch   chan<- *types.LogPackageMeta
	stop atomic.Bool
	done chan struct{}
}

func NewUDPServer(addr string, ch chan<- *types.LogPackageMeta) (*UDPServer, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	udpconn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return nil, err
	}
	return &UDPServer{
		conn: udpconn,
		ch:   ch,
		done: make(chan struct{}),
	}, nil
}

func (srv *UDPServer) Listen() error {
	if srv == nil {
		return nil
	}

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrentLimit)

	buf := make([]byte, bufferSize)

	for {
		if srv.stop.Load() {
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

func (srv *UDPServer) Stop() error {
	if srv == nil {
		return nil
	}
	srv.stop.Store(true)
	<-srv.done
	return srv.conn.Close()
}
