package types

import (
	_types "github.com/504dev/logr-go-client/types"
	"sync"
	"time"
)

type LogPackageJoiner struct {
	sync.RWMutex
	Data map[string]LogPackageRow
	//Queue chan *_types.LogPackage
}

func (m *LogPackageJoiner) Drop(uid string) {
	m.dropSafe(uid)
}
func (m *LogPackageJoiner) dropSafe(uid string) {
	m.Lock()
	delete(m.Data, uid)
	m.Unlock()
}
func (m *LogPackageJoiner) createSafe(lp *_types.LogPackage, lifetime time.Duration) {
	uid := lp.Chunk.Uid
	size := lp.Chunk.N
	m.Lock()
	m.Data[uid] = make(LogPackageRow, size)
	m.Data[uid][lp.Chunk.I] = lp
	m.Unlock()
	go func() {
		<-time.After(lifetime)
		m.dropSafe(uid)
	}()
}

func (m *LogPackageJoiner) addItemSafe(lp *_types.LogPackage) bool {
	m.Lock()
	defer m.Unlock()
	if m.Data[lp.Chunk.Uid] != nil {
		m.Data[lp.Chunk.Uid][lp.Chunk.I] = lp
		return true
	}
	return false
}

func (m *LogPackageJoiner) hasSafe(uid string) bool {
	m.RLock()
	defer m.RUnlock()
	return m.Data[uid] != nil
}

func (m *LogPackageJoiner) completeSafe(uid string) (complete bool, joined *_types.LogPackage) {
	m.RLock()
	defer m.RUnlock()
	if m.Data[uid] == nil {
		return false, nil
	}
	return m.Data[uid].Joined()

}

func (m *LogPackageJoiner) Add(lp *_types.LogPackage, tries int) (complete bool, joined *_types.LogPackage) {
	if lp.Chunk.I == 0 {
		m.createSafe(lp, time.Duration(tries)*time.Second)
		return false, nil
	}

	if !m.hasSafe(lp.Chunk.Uid) {
		if tries > 0 {
			<-time.After(time.Second / time.Duration(tries+1))
			return m.Add(lp, tries-1)
		}
		return false, nil
	}

	m.addItemSafe(lp)

	complete, joined = m.completeSafe(lp.Chunk.Uid)
	if complete {
		m.dropSafe(lp.Chunk.Uid)
	}

	return
}
