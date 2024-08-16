package tests

import (
	"github.com/504dev/logr/dbs/clickhouse/batcher"
	"github.com/504dev/logr/dbs/clickhouse/queue"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type BatchProcessor[T any] interface {
	Push(item T) bool
	Run()
	Stop()
}

func createBatcher[T any](size int, interval time.Duration, handler func([]T)) BatchProcessor[T] {
	return batcher.NewBatcher(size, interval, handler)
}

func createQueue[T any](size int, interval time.Duration, handler func([]T)) BatchProcessor[T] {
	return queue.NewQueue(size, interval, handler)
}

func testBatchProcessor(t *testing.T, create func(int, time.Duration, func([]int)) BatchProcessor[int]) {
	t.Run("Boss test", func(t *testing.T) {
		ts := time.Now()
		expected := [][2]int{{50, 0}, {50, 0}, {30, 1}}
		result := make([][2]int, 0)
		processor := create(50, time.Second, func(batch []int) {
			result = append(result, [2]int{len(batch), int(time.Since(ts).Seconds())})
		})
		t.Log(processor)
		fill := func(n int) {
			for i := 0; i < n; i++ {
				go processor.Push(i)
			}
		}
		go func() {
			go fill(60)
			go fill(70)
			go func() {
				time.Sleep(time.Second * 2)
				processor.Stop()
			}()
		}()

		processor.Run()

		t.Log("result", result)
		assert.Equal(t, expected, result, "Unexpected result")
	})

	t.Run("Basic functionality", func(t *testing.T) {
		results := make([][]int, 0)
		processor := create(3, 100*time.Millisecond, func(batch []int) {
			results = append(results, batch)
		})

		go processor.Run()

		for i := 1; i <= 5; i++ {
			processor.Push(i)
		}

		time.Sleep(150 * time.Millisecond)

		processor.Stop()

		assert.Len(t, results, 2, "Expected 2 batches")
		assert.Equal(t, []int{1, 2, 3}, results[0], "First batch incorrect")
		assert.Equal(t, []int{4, 5}, results[1], "Second batch incorrect")
	})

	t.Run("Time limit", func(t *testing.T) {
		results := make([][]int, 0)
		processor := create(10, 50*time.Millisecond, func(batch []int) {
			results = append(results, batch)
		})

		go processor.Run()

		processor.Push(1)
		processor.Push(2)

		time.Sleep(100 * time.Millisecond)

		processor.Stop()

		assert.Len(t, results, 1, "Expected 1 batch")
		assert.Equal(t, []int{1, 2}, results[0], "Batch incorrect")
	})

	t.Run("Concurrent pushes", func(t *testing.T) {
		results := make([][]int, 0)
		processor := create(5, 100*time.Millisecond, func(batch []int) {
			results = append(results, batch)
		})

		go processor.Run()

		var wg sync.WaitGroup
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				processor.Push(i)
			}(i)
		}

		wg.Wait()
		time.Sleep(150 * time.Millisecond)

		processor.Stop()

		totalItems := 0
		for _, batch := range results {
			totalItems += len(batch)
		}

		assert.Equal(t, 20, totalItems, "Expected 20 items in total")
	})
}

func TestBatcher(t *testing.T) {
	testBatchProcessor(t, createBatcher[int])
}

func TestQueue(t *testing.T) {
	testBatchProcessor(t, createQueue[int])
}
