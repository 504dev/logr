package types

import (
	_types "github.com/504dev/logr-go-client/types"
	"sync"
	"time"
)

type LogPackageJoiner struct {
	delay time.Duration
	tries int
	data  map[string]_types.LogPackageChunks
	mx    sync.RWMutex
}

func NewLogPackageJoiner(delay time.Duration, tries int) *LogPackageJoiner {
	return &LogPackageJoiner{
		delay: delay,
		tries: tries,
		data:  make(map[string]_types.LogPackageChunks, 100),
	}
}

func (j *LogPackageJoiner) dropSafe(uid string) {
	j.mx.Lock()
	delete(j.data, uid)
	j.mx.Unlock()
}

func (j *LogPackageJoiner) createSafe(lp *_types.LogPackage, lifetime time.Duration) {
	uid := lp.Chunk.Uid
	size := lp.Chunk.N
	j.mx.Lock()
	j.data[uid] = make(_types.LogPackageChunks, size)
	j.data[uid][lp.Chunk.I] = lp
	j.mx.Unlock()
	go func() {
		<-time.After(lifetime)
		j.dropSafe(uid)
	}()
}

func (j *LogPackageJoiner) addItemSafe(lp *_types.LogPackage) bool {
	j.mx.Lock()
	defer j.mx.Unlock()
	if j.data[lp.Chunk.Uid] != nil {
		j.data[lp.Chunk.Uid][lp.Chunk.I] = lp
		return true
	}
	return false
}

func (j *LogPackageJoiner) hasSafe(uid string) bool {
	j.mx.RLock()
	defer j.mx.RUnlock()
	return j.data[uid] != nil
}

func (j *LogPackageJoiner) completeSafe(uid string) (complete bool, joined *_types.LogPackage) {
	j.mx.RLock()
	defer j.mx.RUnlock()
	if j.data[uid] == nil {
		return false, nil
	}
	return j.data[uid].Joined()

}

func (j *LogPackageJoiner) Add(lp *_types.LogPackage) (complete bool, joined *_types.LogPackage) {
	return j.addNTries(lp, j.tries)
}

func (j *LogPackageJoiner) addNTries(lp *_types.LogPackage, tries int) (complete bool, joined *_types.LogPackage) {
	if lp.Chunk.I == 0 {
		j.createSafe(lp, time.Duration(j.tries)*j.delay)
		return false, nil
	}

	if !j.hasSafe(lp.Chunk.Uid) {
		if tries > 0 {
			<-time.After(j.delay / time.Duration(tries+1))
			return j.addNTries(lp, tries-1)
		}
		return false, nil
	}

	j.addItemSafe(lp)

	complete, joined = j.completeSafe(lp.Chunk.Uid)
	if complete {
		j.dropSafe(lp.Chunk.Uid)
	}

	return
}
