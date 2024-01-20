package server

import (
	"encoding/json"
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	countModel "github.com/504dev/logr/models/count"
	"github.com/504dev/logr/models/dashkey"
	logModel "github.com/504dev/logr/models/log"
	"github.com/504dev/logr/models/ws"
	"net"
)

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
			Logger.Error("UDP read error: %v", err)
			continue
		}

		Logger.Inc("udp", 1)

		lp := _types.LogPackage{}
		err = json.Unmarshal(buf[0:n], &lp)

		if err != nil {
			Logger.Error("UDP parse json error: %v\n%v", err, string(buf[0:n]))
			continue
		}
		dk, err := dashkey.GetByPubCached(lp.PublicKey)
		if err != nil {
			Logger.Error("UDP dash error: %v", err)
			continue
		}
		if dk == nil {
			Logger.Error("UDP unknown dash pub=%v", lp.PublicKey)
			continue
		}

		if !config.Get().AllowNoCipher {
			lp.Log = nil
			lp.Count = nil
		}

		if lp.CipherLog != "" {
			Logger.Inc("udp:l", 1)
			err = lp.DecryptLog(dk.PrivateKey)
			if err != nil {
				Logger.Error("UDP decrypt log error: %v", err)
			}
		}
		if lp.Log != nil {
			lp.Log.DashId = dk.DashId
			ws.SockMap.PushLog(lp.Log)
			err = logModel.PushToQueue(lp.Log)
			if err != nil {
				Logger.Error("UDP create log error: %v", err)
			}
		}

		if lp.CipherCount != "" {
			Logger.Inc("udp:c", 1)
			err = lp.DecryptCount(dk.PrivateKey)
			if err != nil {
				Logger.Error("UDP decrypt count error: %v", err)
			}
		}
		if lp.Count != nil {
			lp.Count.DashId = dk.DashId
			//Logger.Debug("UDP %v", lp.Count)
			err = countModel.PushToQueue(lp.Count)
			if err != nil {
				Logger.Error("UDP create count error: %v", err)
			}
		}
	}
}
