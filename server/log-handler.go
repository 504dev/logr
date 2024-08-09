package server

import (
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	countModel "github.com/504dev/logr/models/count"
	"github.com/504dev/logr/models/dashkey"
	logModel "github.com/504dev/logr/models/log"
	"github.com/504dev/logr/models/ws"
)

func (srv *LogServer) handleLog(meta *LogPackageMeta) {
	lp := meta.LogPackage
	dk, err := dashkey.GetByPubCached(lp.PublicKey)
	if err != nil {
		Logger.Error("(%v) dash error: %v", meta.Protocol, err)
		return
	}
	if dk == nil {
		Logger.Inc(meta.Protocol+":unknown", 1)
		Logger.Warn("(%v) unknown dash pub=%v, log=%v", meta.Protocol, lp.PublicKey, lp.Log)
		return
	}

	if !config.Get().AllowNoCipher {
		lp.Log = nil
		lp.PlainLog = nil
		lp.Count = nil
	}

	// Handle logs
	if lp.CipherLog != nil || lp.PlainLog != nil || lp.Log != nil {
		Logger.Inc(meta.Protocol+":logs", 1)
		Logger.Inc(meta.Protocol+":logs:bytes", float64(meta.Size))
		go func() {
			if lp.CipherLog != nil || lp.PlainLog != nil {
				if lp.Chunk != nil {
					sig, err := lp.Chunk.CalcSig(dk.PrivateKey)
					if err != nil || lp.Sig != sig {
						Logger.Error("(%v) signature error: %v, %v", meta.Protocol, err, lp.Sig != sig)
						return
					}

					// merging of long logs consisting of several chunks
					if lp.Chunk.N > 1 {
						complete, joined := srv.joiner.Add(lp)
						if !complete {
							return
						}
						lp = joined
					}
				}

				if lp.CipherLog != nil {
					err = lp.DecryptLog(dk.PrivateKey)
					if err != nil {
						Logger.Error("(%v) decrypt log error: %v", meta.Protocol, err)
						return
					}
				} else {
					err = lp.DeserializeLog()
					if err != nil {
						Logger.Error("(%v) deserialize log error: %v", meta.Protocol, err)
						return
					}
				}
			}

			if lp.Log != nil {
				lp.Log.DashId = dk.DashId
				ws.GetSockMap().Push(lp.Log)
				err = logModel.PushToQueue(lp.Log) // TODO srv.storage
				if err != nil {
					Logger.Error("(%v) create log error: %v", meta.Protocol, err)
					return
				}
			}
		}()
		return
	}

	// Handle counts
	if lp.CipherCount != nil || lp.Count != nil {
		Logger.Inc(meta.Protocol+":counts", 1)
		Logger.Inc(meta.Protocol+":counts:bytes", float64(meta.Size))
		go func() {
			if lp.CipherCount != nil {
				err = lp.DecryptCount(dk.PrivateKey)
				if err != nil {
					Logger.Error("(%v) decrypt count error: %v", meta.Protocol, err)
					return
				}
			}

			if lp.Count != nil {
				lp.Count.DashId = dk.DashId
				//Logger.Debug(PROTOCOL+" %v", lp.Count)
				err = countModel.PushToQueue(lp.Count) // TODO srv.storage
				if err != nil {
					Logger.Error("(%v) create count error: %v", meta.Protocol, err)
					return
				}
			}
		}()
	}
}
