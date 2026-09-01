package id_allocator

import (
	"testing"
)

func TestUint32IdAllocator(t *testing.T) {
	alloc := NewUint32IdAllocator(1)
	id1 := alloc.Get()
	id2 := alloc.Get()
	if id1 == id2 {
		t.Error("IdAllocator should generate unique IDs")
	}
}

func TestGenerateSessionId(t *testing.T) {
	id, err := GenerateSessionId()
	if err != nil {
		t.Errorf("GenerateSessionId error: %v", err)
	}
	if len(id) != 32 {
		t.Errorf("SessionId length = %d, want 32", len(id))
	}
}
