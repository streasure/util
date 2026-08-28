package gopool

import (
	"sync/atomic"
	"testing"
)

func TestPool(t *testing.T) {
	var count atomic.Int32
	pool := NewPool(3)

	for i := 0; i < 10; i++ {
		pool.Add(Worker(func() {
			count.Add(1)
		}))
	}

	pool.Stop()

	if count.Load() != 10 {
		t.Errorf("Pool processed %d items, want 10", count.Load())
	}
}

func TestTypedPool(t *testing.T) {
	pool := NewTypedPool[int](3)
	results := make([]<-chan int, 0, 5)

	for i := 0; i < 5; i++ {
		v := i
		ch := pool.Add(func() int {
			return v * 2
		})
		results = append(results, ch)
	}

	pool.Stop()

	for i, ch := range results {
		got := <-ch
		want := i * 2
		if got != want {
			t.Errorf("TypedPool result[%d] = %d, want %d", i, got, want)
		}
	}
}
