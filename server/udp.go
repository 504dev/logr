package server

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	jsoniter "github.com/json-iterator/go"
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

	for {
		buf := make([]byte, 65536)
		size, _, err := pc.ReadFromUDP(buf)

		if err != nil {
			Logger.Error("UDP read error: %v", err)
			continue
		}

		//fmt.Println("DEBUG buf:", string(buf))

		lp := _types.LogPackage{}
		err = jsoniter.Unmarshal(buf[0:size], &lp)

		if err != nil {
			Logger.Error("UDP parse json error: %v\n%v", err, string(buf[0:size]))
			continue
		}

		Handle(&lp, "udp", size)
	}
}
