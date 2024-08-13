package tests

import (
	"github.com/504dev/logr/dbs/clickhouse/batcher"
	"github.com/go-playground/assert/v2"
	"testing"
	"time"
)

func TestBatcher(t *testing.T) {
	ts := time.Now()
	expected := [][2]int{{50, 0}, {50, 0}, {30, 1}}
	result := make([][2]int, 0)
	b := batcher.NewBatcher(50, time.Second, func(batch []int) {
		result = append(result, [2]int{len(batch), int(time.Since(ts).Seconds())})
		t.Log(result)
	})
	t.Log(b)
	fill := func(n int) {
		for i := 0; i < n; i++ {
			go b.Push(i)
		}
	}
	go func() {
		go fill(60)
		go fill(70)
		go func() {
			time.Sleep(time.Second * 2)
			b.Stop()
		}()
	}()

	b.Run()

	t.Log("result", result)
	assert.Equal(t, compareSlices(result, expected), true)
}

func compareSlices(a, b [][2]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
