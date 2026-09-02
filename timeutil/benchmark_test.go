package timeutil

import (
	"testing"
	"time"
)

func BenchmarkDiffNatureDays(b *testing.B) {
	t1 := int64(1609459200)
	t2 := int64(1609545600)
	for i := 0; i < b.N; i++ {
		DiffNatureDays(t1, t2)
	}
}

func BenchmarkDiffDays(b *testing.B) {
	t1 := time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local)
	t2 := time.Date(2021, 1, 3, 18, 0, 0, 0, time.Local)
	for i := 0; i < b.N; i++ {
		DiffDays(t2, t1)
	}
}

func BenchmarkIsSameDay(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		IsSameDay(now, now)
	}
}

func BenchmarkIsSameDayUnix(b *testing.B) {
	ts1 := int64(1609459200)
	ts2 := int64(1609462800)
	for i := 0; i < b.N; i++ {
		IsSameDayUnix(ts1, ts2)
	}
}

func BenchmarkIsSameWeek(b *testing.B) {
	ts1 := int64(1609459200)
	ts2 := int64(1609545600)
	for i := 0; i < b.N; i++ {
		IsSameWeek(ts1, ts2)
	}
}

func BenchmarkIsSameMonth(b *testing.B) {
	ts1 := int64(1609459200)
	ts2 := int64(1609545600)
	for i := 0; i < b.N; i++ {
		IsSameMonth(ts1, ts2)
	}
}

func BenchmarkZeroTimeOfDay(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		ZeroTimeOfDay(now)
	}
}

func BenchmarkGetZeroTime(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		GetZeroTime(now)
	}
}

func BenchmarkNormalizeTimeOfDay(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		NormalizeTimeOfDay(now, 5)
	}
}

func BenchmarkGetTomorrowStamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTomorrowStamp()
	}
}

func BenchmarkTimeStampToString(b *testing.B) {
	ts := int64(1609459200)
	for i := 0; i < b.N; i++ {
		TimeStampToString(ts)
	}
}

func BenchmarkTimeToString(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		TimeToString(now)
	}
}

func BenchmarkGetDateKey(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		GetDateKey(now)
	}
}

func BenchmarkNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Now()
	}
}

func BenchmarkTimestamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Timestamp()
	}
}
