package bucket

import (
	"testing"
	"time"
)

func TestNewBucket(t *testing.T) {
	b := NewBucket(100*time.Millisecond, 10, 2)
	if b == nil {
		t.Fatal("NewBucket returned nil")
	}
	if b.capacity != 10 {
		t.Fatalf("capacity: got %d, want 10", b.capacity)
	}
	if b.availableTokens != 10 {
		t.Fatalf("availableTokens: got %d, want 10", b.availableTokens)
	}
}

func TestTakeAvailable(t *testing.T) {
	b := NewBucket(100*time.Millisecond, 10, 2)

	taken := b.TakeAvailable(time.Now(), 5)
	if taken != 5 {
		t.Fatalf("TakeAvailable: got %d, want 5", taken)
	}

	taken = b.TakeAvailable(time.Now(), 5)
	if taken != 5 {
		t.Fatalf("TakeAvailable second: got %d, want 5", taken)
	}

	taken = b.TakeAvailable(time.Now(), 1)
	if taken != 0 {
		t.Fatalf("TakeAvailable empty: got %d, want 0", taken)
	}
}

func TestTakeAvailableNoLock(t *testing.T) {
	b := NewBucket(100*time.Millisecond, 10, 2)

	taken := b.TakeAvailableNoLock(time.Now(), 3)
	if taken != 3 {
		t.Fatalf("TakeAvailableNoLock: got %d, want 3", taken)
	}

	taken = b.TakeAvailableNoLock(time.Now(), 20)
	if taken != 7 {
		t.Fatalf("TakeAvailableNoLock over: got %d, want 7", taken)
	}
}

func TestTakeAvailableRefill(t *testing.T) {
	b := NewBucket(10*time.Millisecond, 10, 2)

	b.TakeAvailable(time.Now(), 10)

	// Wait for refill: 50ms / 10ms = 5 ticks, 5 * 2 = 10 tokens
	time.Sleep(55 * time.Millisecond)

	taken := b.TakeAvailable(time.Now(), 5)
	if taken != 5 {
		t.Fatalf("TakeAvailable refill: got %d, want 5", taken)
	}
}

func TestTakeAvailableZeroOrNegative(t *testing.T) {
	b := NewBucket(100*time.Millisecond, 10, 2)

	taken := b.TakeAvailable(time.Now(), 0)
	if taken != 0 {
		t.Fatalf("TakeAvailable zero: got %d, want 0", taken)
	}

	taken = b.TakeAvailable(time.Now(), -1)
	if taken != 0 {
		t.Fatalf("TakeAvailable negative: got %d, want 0", taken)
	}
}

func TestTakeAvailableMoreThanCapacity(t *testing.T) {
	b := NewBucket(100*time.Millisecond, 5, 1)

	taken := b.TakeAvailable(time.Now(), 100)
	if taken != 5 {
		t.Fatalf("TakeAvailable over capacity: got %d, want 5", taken)
	}
}
