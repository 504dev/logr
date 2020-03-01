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
	"time"
)

func InfoWs() {
	for {
		time.Sleep(5 * time.Second)
		j, _ := json.MarshalIndent(ws.SockMap, "", "\t")
		logger.Info(string(j))
	}
}

func Init() {
	gin.ForceConsoleColor()

	w, _ := logger.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(os.Stdout, w)

	r := NewRouter()

	r.Run(config.Get().Bind.Http)
}

func Udp() {
	pc, err := net.ListenPacket("udp", config.Get().Bind.Udp)
	if err != nil {
		fmt.Println(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, _, err := pc.ReadFrom(buf)

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

		dash, err := dashboard.GetById(lp.DashId)
		if err != nil {
			fmt.Println("UDP dash error:", err)
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
			err = log.Create(lp.Log)
			if err != nil {
				fmt.Println("UDP create log error", err)
			}
		}

		if lp.Metr != nil {
			fmt.Println(lp.Metr, err)
		}
	}
}
