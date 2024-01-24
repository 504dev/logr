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
func (m *LogPackageJoiner) createSafe(uid string, size int, lifetime time.Duration) {
	m.Lock()
	m.Data[uid] = make(LogPackageRow, size)
	m.Unlock()
	go func() {
		<-time.After(lifetime)
		m.dropSafe(uid)
	}()
}

func (m *LogPackageJoiner) addSafe(lp *_types.LogPackage) {
	m.Lock()
	m.Data[lp.Chunk.Uid][lp.Chunk.I] = lp
	m.Unlock()
}

func (m *LogPackageJoiner) getSafe(uid string, i int) *_types.LogPackage {
	m.RLock()
	defer m.RUnlock()
	if m.Data[uid] == nil {
		return nil
	}
	return m.Data[uid][i]
}

func (m *LogPackageJoiner) hasSafe(uid string) bool {
	m.RLock()
	defer m.RUnlock()
	return m.Data[uid] != nil
}

func (m *LogPackageJoiner) completeSafe(uid string) (complete bool, joined *_types.LogPackage) {
	m.RLock()
	defer m.RUnlock()
	return m.Data[uid].Joined()

}

func (m *LogPackageJoiner) Add(lp *_types.LogPackage, tries int) (complete bool, joined *_types.LogPackage) {
	if lp.Chunk.I == 0 {
		m.createSafe(lp.Chunk.Uid, lp.Chunk.N, time.Duration(tries)*time.Second)
	} else if m.hasSafe(lp.Chunk.Uid) == false {
		if tries > 0 {
			<-time.After(time.Second / time.Duration(tries+1))
			return m.Add(lp, tries-1)
		}
		return false, nil
	}

	m.addSafe(lp)
	return m.completeSafe(lp.Chunk.Uid)
}
