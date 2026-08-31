package util

import (
	"context"
	"testing"
	"time"
)

func TestTimeoutWaitGroupBasic(t *testing.T) {
	wg := NewTimeoutWaitGroup(3)
	if wg == nil {
		t.Fatal("NewTimeoutWaitGroup returned nil")
	}

	err := wg.Add(1)
	if err != nil {
		t.Fatalf("Add error: %v", err)
	}

	err = wg.Add(2)
	if err != nil {
		t.Fatalf("Add error: %v", err)
	}

	wg.Done()
	wg.Done()
	wg.Done()

	finished := wg.WaitTimeout(context.Background(), 100*time.Millisecond)
	if finished {
		t.Fatal("WaitTimeout should not timeout")
	}
}

func TestTimeoutWaitGroupTimeout(t *testing.T) {
	wg := NewTimeoutWaitGroup(1)

	err := wg.Add(1)
	if err != nil {
		t.Fatalf("Add error: %v", err)
	}

	finished := wg.WaitTimeout(context.Background(), 50*time.Millisecond)
	if !finished {
		t.Fatal("WaitTimeout should timeout")
	}
}

func TestTimeoutWaitGroupAddOverflow(t *testing.T) {
	wg := NewTimeoutWaitGroup(1)

	err := wg.Add(2)
	if err != errWaitGroupAddFailed {
		t.Fatalf("expected errWaitGroupAddFailed, got: %v", err)
	}
}

func TestTimeoutWaitGroupNil(t *testing.T) {
	wg := NewTimeoutWaitGroup(0)
	if wg != nil {
		t.Fatal("NewTimeoutWaitGroup(0) should return nil")
	}

	wg = NewTimeoutWaitGroup(-1)
	if wg != nil {
		t.Fatal("NewTimeoutWaitGroup(-1) should return nil")
	}
}

func TestTimeoutWaitGroupZero(t *testing.T) {
	wg := NewTimeoutWaitGroup(1)

	err := wg.Add(0)
	if err != nil {
		t.Fatalf("Add(0) error: %v", err)
	}

	finished := wg.WaitTimeout(context.Background(), 50*time.Millisecond)
	if finished {
		t.Fatal("WaitTimeout should not timeout")
	}
}
