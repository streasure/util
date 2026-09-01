package sys

import (
	"testing"

	"github.com/streasure/util/uuid"
)

func TestNewUUID(t *testing.T) {
	u := uuid.NewUUID()
	if len(u) != 32 {
		t.Errorf("NewUUID length = %d, want 32", len(u))
	}
}

func TestNewUUIDBytes(t *testing.T) {
	b := uuid.NewUUIDBytes()
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
