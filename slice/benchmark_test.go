package slice

import (
	"testing"
)

func BenchmarkUniqueSlice_Small(b *testing.B) {
	s := []int{1, 2, 3, 4, 5, 1, 2, 3}
	for i := 0; i < b.N; i++ {
		UniqueSlice(s)
	}
}

func BenchmarkUniqueSlice_Large(b *testing.B) {
	s := make([]int, 1000)
	for i := range s {
		s[i] = i % 100
	}
	for i := 0; i < b.N; i++ {
		UniqueSlice(s)
	}
}

func BenchmarkUniqueSlice_Strings(b *testing.B) {
	s := []string{"a", "b", "c", "d", "a", "b"}
	for i := 0; i < b.N; i++ {
		UniqueSlice(s)
	}
}

func BenchmarkEqualSlice(b *testing.B) {
	x := []int{1, 2, 3, 4, 5}
	y := []int{5, 4, 3, 2, 1}
	for i := 0; i < b.N; i++ {
		EqualSlice(x, y)
	}
}

func BenchmarkIsNil(b *testing.B) {
	var p *int
	for i := 0; i < b.N; i++ {
		IsNil(p)
	}
}

func BenchmarkIsNil_NonNil(b *testing.B) {
	v := 42
	for i := 0; i < b.N; i++ {
		IsNil(v)
	}
}

func BenchmarkStrToSlice(b *testing.B) {
	s := "hello world"
	for i := 0; i < b.N; i++ {
		StrToSlice(s)
	}
}
