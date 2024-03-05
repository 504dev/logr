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
	"github.com/504dev/logr/types"
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

	var joiner = types.LogPackageJoiner{
		Data: map[string]types.LogPackageRow{},
	}

	for {
		buf := make([]byte, 65536)
		n, _, err := pc.ReadFromUDP(buf)

		if err != nil {
			Logger.Error("UDP read error: %v", err)
			continue
		}

		//fmt.Println("DEBUG buf:", string(buf))

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
			Logger.Inc("udp:unknown", 1)
			Logger.Warn("UDP unknown dash pub=%v, log=%v", lp.PublicKey, lp.Log)
			continue
		}

		if !config.Get().AllowNoCipher {
			lp.Log = nil
			lp.PlainLog = nil
			lp.Count = nil
		}

		// Handle logs
		if lp.CipherLog != nil || lp.PlainLog != nil || lp.Log != nil {
			Logger.Inc("udp:logs", 1)
			Logger.Inc("udp:logs:bytes", float64(n))
			go func() {
				if lp.CipherLog != nil || lp.PlainLog != nil {
					if lp.Chunk != nil {
						sig, err := lp.Chunk.CalcSig(dk.PrivateKey)
						if err != nil || lp.Sig != sig {
							Logger.Error("UDP signature error: %v, %v", err, lp.Sig != sig)
							return
						}

						if lp.Chunk.N > 1 {
							complete, joined := joiner.Add(&lp, 5)
							if !complete {
								return
							}
							joiner.Drop(lp.Chunk.Uid)
							lp = *joined
						}
					}

					if lp.CipherLog != nil {
						err = lp.DecryptLog(dk.PrivateKey)
						if err != nil {
							Logger.Error("UDP decrypt log error: %v", err)
							return
						}
					} else {
						err = lp.DeserializeLog()
						if err != nil {
							Logger.Error("UDP deserialize log error: %v", err)
							return
						}
					}
				}

				if lp.Log != nil {
					lp.Log.DashId = dk.DashId
					ws.SockMap.PushLog(lp.Log)
					err = logModel.PushToQueue(lp.Log)
					if err != nil {
						Logger.Error("UDP create log error: %v", err)
						return
					}
				}
			}()
			continue
		}

		// Handle counts
		if lp.CipherCount != nil || lp.Count != nil {
			Logger.Inc("udp:counts", 1)
			Logger.Inc("udp:counts:bytes", float64(n))
			go func() {
				if lp.CipherCount != nil {
					err = lp.DecryptCount(dk.PrivateKey)
					if err != nil {
						Logger.Error("UDP decrypt count error: %v", err)
						return
					}
				}

				if lp.Count != nil {
					lp.Count.DashId = dk.DashId
					//Logger.Debug("UDP %v", lp.Count)
					err = countModel.PushToQueue(lp.Count)
					if err != nil {
						Logger.Error("UDP create count error: %v", err)
						return
					}
				}
			}()
		}
	}
}
