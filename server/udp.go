package server

import (
	"encoding/json"
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	"net"
)

func MustListenUDP() {
	if err := ListenUDP(); err != nil {
		panic(err)
	}
}
func ListenUDP() error {
	serverAddr, err := net.ResolveUDPAddr("udp", config.Get().Bind.Udp)
	if err != nil {
		return err
	}
	pc, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}
	defer pc.Close()

	buf := make([]byte, 65536)

	for {
		size, _, err := pc.ReadFromUDP(buf)
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

			Handle(&lp, "udp", size)
		}()
	}
}
