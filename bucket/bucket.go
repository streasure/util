package bucket

import (
	"sync"
	"time"
)

type Bucket struct {
	startTime       time.Time
	capacity        int64
	quantum         int64
	fillInterval    time.Duration
	availableTokens int64
	latestTick      int64
	mu              sync.Mutex
}

func NewBucket(fillInterval time.Duration, capacity, quantum int64) *Bucket {
	return &Bucket{
		startTime:       time.Now(),
		fillInterval:    fillInterval,
		capacity:        capacity,
		quantum:         quantum,
		availableTokens: capacity,
	}
}

func (tb *Bucket) TakeAvailableNoLock(now time.Time, count int64) int64 {
	if count <= 0 {
		return 0
	}

	tick := tb.currentTick(now)
	tb.adjustAvailableTokens(tick)
	if tb.availableTokens <= 0 {
		return 0
	}

	if count > tb.availableTokens {
		count = tb.availableTokens
	}
	tb.availableTokens -= count
	return count
}

func (tb *Bucket) TakeAvailable(now time.Time, count int64) int64 {
	if count <= 0 {
		return 0
	}

	tb.mu.Lock()
	defer tb.mu.Unlock()

	tick := tb.currentTick(now)
	tb.adjustAvailableTokens(tick)
	if tb.availableTokens <= 0 {
		return 0
	}

	if count > tb.availableTokens {
		count = tb.availableTokens
	}
	tb.availableTokens -= count
	return count
}

func (tb *Bucket) currentTick(now time.Time) int64 {
	return int64(now.Sub(tb.startTime) / tb.fillInterval)
}

func (tb *Bucket) adjustAvailableTokens(tick int64) {
	if tb.availableTokens >= tb.capacity {
		return
	}

	tb.availableTokens += (tick - tb.latestTick) * tb.quantum
	if tb.availableTokens > tb.capacity {
		tb.availableTokens = tb.capacity
	}
	tb.latestTick = tick
}
