package timeutil

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

func TestTimeToString(t *testing.T) {
	tm := time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local)
	got := TimeToString(tm)
	if got == "" {
		t.Error("TimeToString returned empty string")
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

func TestDiffNatureDays_SameTime(t *testing.T) {
	ts := int64(1609459200)
	got := DiffNatureDays(ts, ts)
	if got != 0 {
		t.Errorf("DiffNatureDays(same) = %d, want 0", got)
	}
}

func TestDiffNatureDays_ReversedOrder(t *testing.T) {
	t1 := int64(1609545600) // 2021-01-02
	t2 := int64(1609459200) // 2021-01-01
	got := DiffNatureDays(t1, t2)
	if got != 1 {
		t.Errorf("DiffNatureDays(reversed) = %d, want 1", got)
	}
}

func TestDiffNatureDays_MultipleDays(t *testing.T) {
	t1 := int64(1609459200) // 2021-01-01
	t2 := int64(1609718400) // 2021-01-04
	got := DiffNatureDays(t1, t2)
	if got != 3 {
		t.Errorf("DiffNatureDays(multiple) = %d, want 3", got)
	}
}

func TestDiffDays(t *testing.T) {
	t1 := time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local)
	t2 := time.Date(2021, 1, 3, 18, 0, 0, 0, time.Local)
	got := DiffDays(t2, t1)
	if got != 2 {
		t.Errorf("DiffDays = %d, want 2", got)
	}
}

func TestZeroTimeOfDay(t *testing.T) {
	now := time.Now()
	zero := ZeroTimeOfDay(now)
	if zero.Hour() != 0 || zero.Minute() != 0 || zero.Second() != 0 {
		t.Error("ZeroTimeOfDay did not return midnight")
	}
}

func TestNormalizeTimeOfDay(t *testing.T) {
	now := time.Date(2021, 1, 5, 3, 0, 0, 0, time.Local)
	result := NormalizeTimeOfDay(now, 5)
	if result.Hour() != 5 {
		t.Errorf("NormalizeTimeOfDay hour = %d, want 5", result.Hour())
	}
	if result.Day() != now.Day()-1 {
		t.Errorf("NormalizeTimeOfDay day = %d, want %d", result.Day(), now.Day()-1)
	}
}

func TestNormalizeTimeOfDay_AfterStartHour(t *testing.T) {
	now := time.Date(2021, 1, 1, 10, 0, 0, 0, time.Local)
	result := NormalizeTimeOfDay(now, 5)
	if result.Hour() != 5 {
		t.Errorf("NormalizeTimeOfDay hour = %d, want 5", result.Hour())
	}
	if result.Day() != now.Day() {
		t.Error("NormalizeTimeOfDay should stay on same day")
	}
}

func TestGetTomorrowStamp(t *testing.T) {
	tomorrow := GetTomorrowStamp()
	if tomorrow.Hour() != 0 || tomorrow.Minute() != 0 || tomorrow.Second() != 0 {
		t.Error("GetTomorrowStamp should return midnight")
	}
	if tomorrow.Day() != time.Now().Day()+1 {
		t.Error("GetTomorrowStamp should be next day")
	}
}

func TestIsSameDay(t *testing.T) {
	now := time.Now()
	if !IsSameDay(now, now) {
		t.Error("IsSameDay(now, now) should be true")
	}
}

func TestIsSameDay_DifferentDays(t *testing.T) {
	t1 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)
	t2 := time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local)
	if IsSameDay(t1, t2) {
		t.Error("IsSameDay(different days) should be false")
	}
}

func TestIsSameDayUnix(t *testing.T) {
	ts1 := int64(1609459200) // 2021-01-01 00:00:00 UTC
	ts2 := int64(1609462800) // 2021-01-01 01:00:00 UTC
	if !IsSameDayUnix(ts1, ts2) {
		t.Error("IsSameDayUnix(same day) should be true")
	}
}

func TestIsSameWeek(t *testing.T) {
	ts1 := int64(1609459200) // 2021-01-01
	ts2 := int64(1609545600) // 2021-01-02
	if !IsSameWeek(ts1, ts2) {
		t.Error("IsSameWeek(same week) should be true")
	}
}

func TestIsSameMonth(t *testing.T) {
	ts1 := int64(1609459200) // 2021-01-01
	ts2 := int64(1609545600) // 2021-01-02
	if !IsSameMonth(ts1, ts2) {
		t.Error("IsSameMonth(same month) should be true")
	}
}

func TestIsSameMonth_DifferentMonths(t *testing.T) {
	ts1 := int64(1609459200) // 2021-01-01
	ts2 := int64(1612137600) // 2021-02-01
	if IsSameMonth(ts1, ts2) {
		t.Error("IsSameMonth(different months) should be false")
	}
}

func TestIsToday(t *testing.T) {
	now := time.Now()
	if !IsToday(now) {
		t.Error("IsToday(now) should be true")
	}
}

func TestIsTodayUnix(t *testing.T) {
	ts := time.Now().Unix()
	if !IsTodayUnix(ts) {
		t.Error("IsTodayUnix(now) should be true")
	}
}

func TestGetZeroTime(t *testing.T) {
	now := time.Now()
	zero := GetZeroTime(now)
	if zero.Hour() != 0 || zero.Minute() != 0 || zero.Second() != 0 {
		t.Error("GetZeroTime did not return midnight")
	}
}

func TestGetTimeByHour(t *testing.T) {
	now := time.Now()
	result := GetTimeByHour(now, 14)
	if result.Hour() != 14 {
		t.Errorf("GetTimeByHour = %d, want 14", result.Hour())
	}
}

func TestGetDateKey(t *testing.T) {
	now := time.Now()
	key := GetDateKey(now)
	if len(key) != 6 {
		t.Errorf("GetDateKey = %q, want 6 chars", key)
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
