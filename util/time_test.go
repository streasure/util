package util

import (
	"testing"
	"time"
)

func TestTimeStampToString(t *testing.T) {
	ts := int64(1609459200) // 2021-01-01 00:00:00 UTC
	got := TimeStampToString(ts)
	if got == "" {
		t.Error("TimeStampToString returned empty string")
	}
}

func TestDiffNatureDays(t *testing.T) {
	t1 := int64(1609459200) // 2021-01-01
	t2 := int64(1609545600) // 2021-01-02
	got := DiffNatureDays(t1, t2)
	if got != 1 {
		t.Errorf("DiffNatureDays = %d, want 1", got)
	}
}

func TestIsSameDay(t *testing.T) {
	now := time.Now()
	if !IsSameDay(now, now) {
		t.Error("IsSameDay(now, now) should be true")
	}
}

func TestGetZeroTime(t *testing.T) {
	now := time.Now()
	zero := GetZeroTime(now)
	if zero.Hour() != 0 || zero.Minute() != 0 || zero.Second() != 0 {
		t.Error("GetZeroTime did not return midnight")
	}
}

func TestGetDateKey(t *testing.T) {
	now := time.Now()
	key := GetDateKey(now)
	if len(key) != 6 {
		t.Errorf("GetDateKey = %q, want 6 chars", key)
	}
}
