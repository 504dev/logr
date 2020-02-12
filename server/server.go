package server

import (
	"encoding/json"
	"fmt"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/log"
	"github.com/504dev/kidlog/types"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"os"
)

func Init() {
	gin.ForceConsoleColor()

	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger.Logr)

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
			fmt.Println("UDP parse json error:", err)
			continue
		}

		dash, err := dashboard.GetByPub(lp.PublicKey)
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
			//fmt.Println(lp.Log)
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
