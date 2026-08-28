package util

import (
	"testing"
)

func TestNewUUID(t *testing.T) {
	uuid := NewUUID()
	if len(uuid) != 32 {
		t.Errorf("NewUUID length = %d, want 32", len(uuid))
	}
}

func TestNewUUIDBytes(t *testing.T) {
	b := NewUUIDBytes()
	if len(b) != 16 {
		t.Errorf("NewUUIDBytes length = %d, want 16", len(b))
	}
}

func TestGoRoutineId(t *testing.T) {
	id := GoRoutineId()
	if id <= 0 {
		t.Errorf("GoRoutineId = %d, want > 0", id)
	}
}
