package mathutil

import (
	"testing"
)

func BenchmarkClamp_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Clamp(5, 0, 10)
	}
}

func BenchmarkClamp_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Clamp(5.0, 0.0, 10.0)
	}
}

func BenchmarkMax_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Max(5, 10)
	}
}

func BenchmarkMax_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Max(5.0, 10.0)
	}
}

func BenchmarkMin_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Min(5, 10)
	}
}

func BenchmarkMin_Float64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Min(5.0, 10.0)
	}
}

func BenchmarkAbs_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Abs(-5)
	}
}

func BenchmarkAbs_Int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Abs(int64(-5))
	}
}

func BenchmarkOperator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Operator(">", 5, 10)
	}
}

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandInt(1, 100)
	}
}

func BenchmarkRandInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandInt32(1, 100)
	}
}

func BenchmarkRandInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandInt64(1, 100)
	}
}

func BenchmarkClampFloat64Math(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ClampFloat64Math(5.0, 0.0, 10.0)
	}
}
