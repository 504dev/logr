package types

import (
	_types "github.com/504dev/logr-go-client/types"
	"sync"
	"time"
)

type LogPackageJoiner struct {
	mu    sync.RWMutex
	delay time.Duration
	tries int
	data  map[string]_types.LogPackageChunks
}

func NewLogPackageJoiner(delay time.Duration, tries int) *LogPackageJoiner {
	return &LogPackageJoiner{
		delay: delay,
		tries: tries,
		data:  make(map[string]_types.LogPackageChunks, 100),
	}
}

func (j *LogPackageJoiner) dropSafe(uid string) {
	j.mu.Lock()
	delete(j.data, uid)
	j.mu.Unlock()
}

func (j *LogPackageJoiner) createSafe(lp *_types.LogPackage, lifetime time.Duration) {
	uid := lp.Chunk.Uid
	size := lp.Chunk.N

	j.mu.Lock()
	j.data[uid] = make(_types.LogPackageChunks, size)
	j.data[uid][lp.Chunk.I] = lp
	j.mu.Unlock()

	go func() {
		<-time.After(lifetime)
		j.dropSafe(uid)
	}()
}

func (j *LogPackageJoiner) addItemSafe(lp *_types.LogPackage) bool {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.data[lp.Chunk.Uid] != nil {
		j.data[lp.Chunk.Uid][lp.Chunk.I] = lp
		return true
	}

	return false
}

func (j *LogPackageJoiner) hasSafe(uid string) bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.data[uid] != nil
}

func (j *LogPackageJoiner) completeSafe(uid string) (bool, *_types.LogPackage) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	if j.data[uid] == nil {
		return false, nil
	}
	return j.data[uid].Joined()
}

func (j *LogPackageJoiner) Add(lp *_types.LogPackage) (bool, *_types.LogPackage) {
	return j.addNTries(lp, j.tries)
}

func (j *LogPackageJoiner) addNTries(lp *_types.LogPackage, tries int) (bool, *_types.LogPackage) {
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

	complete, joined := j.completeSafe(lp.Chunk.Uid)
	if complete {
		j.dropSafe(lp.Chunk.Uid)
	}

	return complete, joined
}
