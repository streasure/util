package mathx

import (
	"testing"
)

func BenchmarkV2_Len(b *testing.B) {
	v := NewV2(3.0, 4.0)
	for i := 0; i < b.N; i++ {
		v.Len()
	}
}

func BenchmarkV2_LenSqrt(b *testing.B) {
	v := NewV2(3.0, 4.0)
	for i := 0; i < b.N; i++ {
		v.LenSqrt()
	}
}

func BenchmarkV2_Normalize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := NewV2(3.0, 4.0)
		v.Normalize()
	}
}

func BenchmarkV2_Add(b *testing.B) {
	v1 := NewV2(1.0, 2.0)
	v2 := NewV2(3.0, 4.0)
	for i := 0; i < b.N; i++ {
		v1.Add(v2)
	}
}

func BenchmarkV2_Sub(b *testing.B) {
	v1 := NewV2(1.0, 2.0)
	v2 := NewV2(3.0, 4.0)
	for i := 0; i < b.N; i++ {
		v1.Sub(v2)
	}
}

func BenchmarkV2_Mul(b *testing.B) {
	v := NewV2(1.0, 2.0)
	for i := 0; i < b.N; i++ {
		v.Mul(2.0)
	}
}

func BenchmarkV2_Div(b *testing.B) {
	v := NewV2(1.0, 2.0)
	for i := 0; i < b.N; i++ {
		v.Div(2.0)
	}
}

func BenchmarkV2_Dot(b *testing.B) {
	v1 := NewV2(1.0, 2.0)
	v2 := NewV2(3.0, 4.0)
	for i := 0; i < b.N; i++ {
		v1.Dot(&v2)
	}
}

func BenchmarkAdd(b *testing.B) {
	v1 := NewV2(1.0, 2.0)
	v2 := NewV2(3.0, 4.0)
	for i := 0; i < b.N; i++ {
		Add(v1, v2)
	}
}

func BenchmarkSub(b *testing.B) {
	v1 := NewV2(1.0, 2.0)
	v2 := NewV2(3.0, 4.0)
	for i := 0; i < b.N; i++ {
		Sub(v1, v2)
	}
}

func BenchmarkMul(b *testing.B) {
	v := NewV2(1.0, 2.0)
	for i := 0; i < b.N; i++ {
		Mul(v, 2.0)
	}
}

func BenchmarkDiv(b *testing.B) {
	v := NewV2(1.0, 2.0)
	for i := 0; i < b.N; i++ {
		Div(v, 2.0)
	}
}

func BenchmarkDot(b *testing.B) {
	v1 := NewV2(1.0, 2.0)
	v2 := NewV2(3.0, 4.0)
	for i := 0; i < b.N; i++ {
		Dot(v1, v2)
	}
}

func BenchmarkIsFloatSame(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsFloatSame(1.0, 1.0)
	}
}
