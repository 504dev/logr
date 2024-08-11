package server

import (
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/dashkey"
	"github.com/504dev/logr/types"
)

func (srv *LogServer) handle(meta *types.LogPackageMeta) {
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
				srv.sockmap.Push(lp.Log)
				err = srv.logStorage.Store(lp.Log)
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
				err = srv.countStorage.Store(lp.Count)
				if err != nil {
					Logger.Error("(%v) create count error: %v", meta.Protocol, err)
					return
				}
			}
		}()
	}
}
