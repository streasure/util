package now

import (
	"testing"
	"time"

	gametime "github.com/streasure/util/gametime"
)

func BenchmarkBeginningOfDay(b *testing.B) {
	gametime.SetOffset(0)
	for i := 0; i < b.N; i++ {
		BeginningOfDay()
	}
}

func BenchmarkBeginningOfWeek(b *testing.B) {
	gametime.SetOffset(0)
	for i := 0; i < b.N; i++ {
		BeginningOfWeek()
	}
}

func BenchmarkBeginningOfMonth(b *testing.B) {
	gametime.SetOffset(0)
	for i := 0; i < b.N; i++ {
		BeginningOfMonth()
	}
}

func BenchmarkBeginningOfYear(b *testing.B) {
	gametime.SetOffset(0)
	for i := 0; i < b.N; i++ {
		BeginningOfYear()
	}
}

func BenchmarkEndOfDay(b *testing.B) {
	gametime.SetOffset(0)
	for i := 0; i < b.N; i++ {
		EndOfDay()
	}
}

func BenchmarkEndOfWeek(b *testing.B) {
	gametime.SetOffset(0)
	for i := 0; i < b.N; i++ {
		EndOfWeek()
	}
}

func BenchmarkEndOfMonth(b *testing.B) {
	gametime.SetOffset(0)
	for i := 0; i < b.N; i++ {
		EndOfMonth()
	}
}

func BenchmarkEndOfYear(b *testing.B) {
	gametime.SetOffset(0)
	for i := 0; i < b.N; i++ {
		EndOfYear()
	}
}

func BenchmarkWith(b *testing.B) {
	tm := time.Date(2021, 1, 15, 12, 30, 0, 0, time.Local)
	for i := 0; i < b.N; i++ {
		With(tm)
	}
}

func BenchmarkNew(b *testing.B) {
	tm := time.Date(2021, 1, 15, 12, 30, 0, 0, time.Local)
	for i := 0; i < b.N; i++ {
		New(tm)
	}
}
