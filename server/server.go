package server

import (
	"encoding/json"
	"fmt"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/log"
	"github.com/504dev/kidlog/models/ws"
	"github.com/504dev/kidlog/types"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"os"
)

func ListenHTTP() error {
	gin.ForceConsoleColor()

	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger.Gin)

	r := NewRouter()

	return r.Run(config.Get().Bind.Http)
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
		n, _, err := pc.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("UDP read error:", err)
			continue
		}

		lp := types.LogPackage{}
		err = json.Unmarshal(buf[0:n], &lp)
		if err != nil {
			fmt.Println("UDP parse json error:", err, string(buf[0:n]))
			continue
		}
		var dash *types.Dashboard
		if lp.DashId != 0 {
			dash, err = dashboard.GetByIdCached(lp.DashId)
		} else {
			dash, err = dashboard.GetByPubCached(lp.PublicKey)
		}
		if err != nil {
			fmt.Println("UDP dash error:", err)
			continue
		}
		if dash == nil {
			fmt.Println("UDP unknown dash")
			continue
		}
		err = lp.DecryptLog(dash.PrivateKey)
		if err != nil {
			fmt.Println("UDP decrypt error:", err)
			continue
		}

		if lp.Log != nil {
			lp.Log.DashId = dash.Id
			//fmt.Println(lp.Log)
			ws.SockMap.PushLog(lp.Log)
			err = log.PushToQueue(lp.Log)
			if err != nil {
				fmt.Println("UDP create log error", err)
			}
		}

		if lp.Counter != nil {
			fmt.Println(lp.Counter, err)
		}
	}
}
