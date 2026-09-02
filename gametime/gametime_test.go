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

func TestGetOffset(t *testing.T) {
	orig := GetOffset()
	SetOffset(time.Hour)
	if GetOffset() != time.Hour {
		t.Error("GetOffset failed")
	}
	SetOffset(orig)
}

func TestUnix(t *testing.T) {
	ts := int64(1609459200)
	tm := Unix(ts, 0)
	if tm.Unix() != ts {
		t.Errorf("Unix = %v, want %v", tm.Unix(), ts)
	}
}

func TestTimestamp(t *testing.T) {
	ts := Timestamp()
	expected := Now().Unix()
	if ts != expected {
		t.Errorf("Timestamp = %d, want %d", ts, expected)
	}
}

func TestSince(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	d := Since(past)
	if d < time.Hour-time.Second || d > time.Hour+time.Second {
		t.Errorf("Since = %v, want ~1 hour", d)
	}
}

func TestUntil(t *testing.T) {
	future := time.Now().Add(time.Hour)
	d := Until(future)
	if d < time.Hour-time.Second || d > time.Hour+time.Second {
		t.Errorf("Until = %v, want ~1 hour", d)
	}
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

func TestRefTimeNextNDayResetTime(t *testing.T) {
	rt := NewRefTime(DailyTime{5, 0, 0})
	now := time.Now()
	resetTime := rt.NextNDayResetTime(now, 1)
	if resetTime.Hour() != 5 {
		t.Errorf("NextNDayResetTime hour = %d, want 5", resetTime.Hour())
	}
}

func TestRefTimeIsSameWeek(t *testing.T) {
	rt := NewRefTime(DailyTime{0, 0, 0})
	now := time.Now()
	if !rt.IsSameWeek(now, now) {
		t.Error("IsSameWeek(now, now) should be true")
	}
}

func TestRefTimeIsSameMonth(t *testing.T) {
	rt := NewRefTime(DailyTime{0, 0, 0})
	now := time.Now()
	if !rt.IsSameMonth(now, now) {
		t.Error("IsSameMonth(now, now) should be true")
	}
}
