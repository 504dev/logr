package tests

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/types"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Предположим, что _types.LogPackage и _types.LogPackageChunks определены в отдельном пакете
// Для целей тестирования мы создадим упрощенные версии этих структур

func TestLogPackageJoiner_Add(t *testing.T) {
	joiner := types.NewLogPackageJoiner(time.Millisecond*10, 3)

	// Тест добавления первого чанка
	lp1 := &_types.LogPackage{Chunk: &_types.ChunkInfo{Uid: "test1", I: 0, N: 3}, PlainLog: []byte("chunk1")}
	complete, joined := joiner.Add(lp1)
	assert.False(t, complete)
	assert.Nil(t, joined)

	// Тест добавления второго чанка
	lp2 := &_types.LogPackage{Chunk: &_types.ChunkInfo{Uid: "test1", I: 1, N: 3}, PlainLog: []byte("chunk2")}
	complete, joined = joiner.Add(lp2)
	assert.False(t, complete)
	assert.Nil(t, joined)

	// Тест добавления последнего чанка и завершения пакета
	lp3 := &_types.LogPackage{Chunk: &_types.ChunkInfo{Uid: "test1", I: 2, N: 3}, PlainLog: []byte("chunk3")}
	complete, joined = joiner.Add(lp3)
	assert.True(t, complete)
	assert.NotNil(t, joined)
	assert.Equal(t, []byte("chunk1chunk2chunk3"), joined.PlainLog)
}

func TestLogPackageJoiner_AddOutOfOrder(t *testing.T) {
	joiner := types.NewLogPackageJoiner(time.Millisecond*10, 3)

	lp1 := &_types.LogPackage{Chunk: &_types.ChunkInfo{Uid: "test2", I: 0, N: 3}, PlainLog: []byte("chunk1")}
	lp2 := &_types.LogPackage{Chunk: &_types.ChunkInfo{Uid: "test2", I: 1, N: 3}, PlainLog: []byte("chunk2")}
	lp3 := &_types.LogPackage{Chunk: &_types.ChunkInfo{Uid: "test2", I: 2, N: 3}, PlainLog: []byte("chunk3")}
	chunks := _types.LogPackageChunks{lp2, lp1, lp3}

	var wg sync.WaitGroup
	var joined *_types.LogPackage

	wg.Add(len(chunks))
	for i, lp := range chunks {
		go func(i int, lp *_types.LogPackage) {
			defer wg.Done()
			<-time.After(time.Millisecond * time.Duration(i+1))
			complete, j := joiner.Add(lp)
			if complete {
				joined = j
				assert.NotNil(t, j)
			} else {
				assert.Nil(t, j)
			}
		}(i, lp)
	}

	wg.Wait()

	assert.Equal(t, []byte("chunk1chunk2chunk3"), joined.PlainLog)
}

func TestLogPackageJoiner_AddTimeout(t *testing.T) {
	joiner := types.NewLogPackageJoiner(time.Millisecond*10, 2)

	lp1 := &_types.LogPackage{Chunk: &_types.ChunkInfo{Uid: "test3", I: 0, N: 2}, PlainLog: []byte("chunk1")}
	joiner.Add(lp1)

	time.Sleep(time.Millisecond * 30) // Ждем, пока пакет устареет

	lp2 := &_types.LogPackage{Chunk: &_types.ChunkInfo{Uid: "test3", I: 1, N: 2}, PlainLog: []byte("chunk2")}
	complete, joined := joiner.Add(lp2)
	assert.False(t, complete)
	assert.Nil(t, joined)
}

func TestLogPackageJoiner_AddCiphered(t *testing.T) {
	joiner := types.NewLogPackageJoiner(time.Millisecond*10, 3)

	lp1 := &_types.LogPackage{Chunk: &_types.ChunkInfo{Uid: "test4", I: 0, N: 2}, CipherLog: []byte("encrypted1")}
	joiner.Add(lp1)

	lp2 := &_types.LogPackage{Chunk: &_types.ChunkInfo{Uid: "test4", I: 1, N: 2}, CipherLog: []byte("encrypted2")}
	complete, joined := joiner.Add(lp2)
	assert.True(t, complete)
	assert.NotNil(t, joined)
	assert.Equal(t, []byte("encrypted1encrypted2"), joined.CipherLog)
}
