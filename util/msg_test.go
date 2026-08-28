package util

import (
	"testing"
	"time"
)

func TestSyncMessage(t *testing.T) {
	msg := NewSyncMessage()
	done := make(chan struct{})
	go func() {
		time.Sleep(10 * time.Millisecond)
		msg.Wait()
		close(done)
	}()
	time.Sleep(20 * time.Millisecond)
	msg.Done()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("SyncMessage.Wait should return after Done")
	}
}

func TestAsyncMessage(t *testing.T) {
	msg := NewAsyncMessage()
	msg.Done()
	msg.Wait()
}
