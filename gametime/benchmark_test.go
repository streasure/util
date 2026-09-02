package gametime

import (
	"testing"
	"time"
)

func BenchmarkNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Now()
	}
}

func BenchmarkSetOffset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetOffset(time.Duration(i))
	}
}

func BenchmarkGetOffset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetOffset()
	}
}

func BenchmarkTimestamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Timestamp()
	}
}

func BenchmarkSince(b *testing.B) {
	past := time.Now().Add(-time.Hour)
	for i := 0; i < b.N; i++ {
		Since(past)
	}
}

func BenchmarkUntil(b *testing.B) {
	future := time.Now().Add(time.Hour)
	for i := 0; i < b.N; i++ {
		Until(future)
	}
}

func BenchmarkRefTimeIsSameDay(b *testing.B) {
	rt := NewRefTime(DailyTime{0, 0, 0})
	now := time.Now()
	for i := 0; i < b.N; i++ {
		rt.IsSameDay(now, now)
	}
}

func BenchmarkRefTimeSubDay(b *testing.B) {
	rt := NewRefTime(DailyTime{0, 0, 0})
	t1 := time.Now()
	t2 := t1.Add(48 * time.Hour)
	for i := 0; i < b.N; i++ {
		rt.SubDay(t2, t1)
	}
}

func BenchmarkRefTimeNextNDayResetTime(b *testing.B) {
	rt := NewRefTime(DailyTime{5, 0, 0})
	now := time.Now()
	for i := 0; i < b.N; i++ {
		rt.NextNDayResetTime(now, 1)
	}
}

func BenchmarkRefTimeIsSameWeek(b *testing.B) {
	rt := NewRefTime(DailyTime{0, 0, 0})
	now := time.Now()
	for i := 0; i < b.N; i++ {
		rt.IsSameWeek(now, now)
	}
}

func BenchmarkRefTimeIsSameMonth(b *testing.B) {
	rt := NewRefTime(DailyTime{0, 0, 0})
	now := time.Now()
	for i := 0; i < b.N; i++ {
		rt.IsSameMonth(now, now)
	}
}
