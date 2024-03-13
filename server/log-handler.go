package server

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	countModel "github.com/504dev/logr/models/count"
	"github.com/504dev/logr/models/dashkey"
	logModel "github.com/504dev/logr/models/log"
	"github.com/504dev/logr/models/ws"
	"github.com/504dev/logr/types"
	"strings"
)

var joiner = types.LogPackageJoiner{
	Data: map[string]types.LogPackageRow{},
}

func Handle(lp *_types.LogPackage, protocol string, size int) {
	PROTOCOL := strings.ToUpper(protocol)

	dk, err := dashkey.GetByPubCached(lp.PublicKey)
	if err != nil {
		Logger.Error(PROTOCOL+" dash error: %v", err)
		return
	}
	if dk == nil {
		Logger.Inc(protocol+":unknown", 1)
		Logger.Warn(PROTOCOL+" unknown dash pub=%v, log=%v", lp.PublicKey, lp.Log)
		return
	}

	if !config.Get().AllowNoCipher {
		lp.Log = nil
		lp.PlainLog = nil
		lp.Count = nil
	}

	// Handle logs
	if lp.CipherLog != nil || lp.PlainLog != nil || lp.Log != nil {
		Logger.Inc(protocol+":logs", 1)
		Logger.Inc(protocol+":logs:bytes", float64(size))
		go func() {
			if lp.CipherLog != nil || lp.PlainLog != nil {
				if lp.Chunk != nil {
					sig, err := lp.Chunk.CalcSig(dk.PrivateKey)
					if err != nil || lp.Sig != sig {
						Logger.Error(PROTOCOL+" signature error: %v, %v", err, lp.Sig != sig)
						return
					}

					if lp.Chunk.N > 1 {
						complete, joined := joiner.Add(lp, 5)
						if !complete {
							return
						}
						joiner.Drop(lp.Chunk.Uid)
						lp = joined
					}
				}

				if lp.CipherLog != nil {
					err = lp.DecryptLog(dk.PrivateKey)
					if err != nil {
						Logger.Error(PROTOCOL+" decrypt log error: %v", err)
						return
					}
				} else {
					err = lp.DeserializeLog()
					if err != nil {
						Logger.Error(PROTOCOL+" deserialize log error: %v", err)
						return
					}
				}
			}

			if lp.Log != nil {
				lp.Log.DashId = dk.DashId
				ws.SockMap.PushLog(lp.Log)
				err = logModel.PushToQueue(lp.Log)
				if err != nil {
					Logger.Error(PROTOCOL+" create log error: %v", err)
					return
				}
			}
		}()
		return
	}

	// Handle counts
	if lp.CipherCount != nil || lp.Count != nil {
		Logger.Inc(protocol+":counts", 1)
		Logger.Inc(protocol+":counts:bytes", float64(size))
		go func() {
			if lp.CipherCount != nil {
				err = lp.DecryptCount(dk.PrivateKey)
				if err != nil {
					Logger.Error(PROTOCOL+" decrypt count error: %v", err)
					return
				}
			}

			if lp.Count != nil {
				lp.Count.DashId = dk.DashId
				//Logger.Debug(PROTOCOL+" %v", lp.Count)
				err = countModel.PushToQueue(lp.Count)
				if err != nil {
					Logger.Error(PROTOCOL+" create count error: %v", err)
					return
				}
			}
		}()
	}
}
