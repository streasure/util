package gametime

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	before := time.Now()
	now := Now()
	after := time.Now()
	if now.Before(before.Add(-time.Second)) || now.After(after.Add(time.Second)) {
		t.Error("Now() is not within expected range")
	}
}

func TestSetOffset(t *testing.T) {
	orig := SetOffset(time.Hour)
	defer SetOffset(0)

	now := Now()
	expected := time.Now().Add(time.Hour)
	if now.Sub(expected) > time.Second {
		t.Error("SetOffset failed")
	}
	SetOffset(orig)
}

func TestRefTimeIsSameDay(t *testing.T) {
	rt := NewRefTime(DailyTime{0, 0, 0})
	now := time.Now()
	if !rt.IsSameDay(now, now) {
		t.Error("IsSameDay(now, now) should be true")
	}
}

func TestRefTimeSubDay(t *testing.T) {
	rt := NewRefTime(DailyTime{0, 0, 0})
	t1 := time.Now()
	t2 := t1.Add(48 * time.Hour)
	got := rt.SubDay(t2, t1)
	if got != 2 {
		t.Errorf("SubDay = %d, want 2", got)
	}
}
