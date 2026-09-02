package rand

import (
	"testing"
)

func BenchmarkIntN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntN(100)
	}
}

func BenchmarkInt32N(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int32N(100)
	}
}

func BenchmarkFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64()
	}
}

func BenchmarkRangeInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RangeInt(1, 100)
	}
}

func BenchmarkRangeInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RangeInt32(1, 100)
	}
}

func BenchmarkRangeInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RangeInt64(1, 100)
	}
}

func BenchmarkSliceOne(b *testing.B) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < b.N; i++ {
		SliceOne(arr)
	}
}

func BenchmarkSliceN(b *testing.B) {
	arr := make([]int, 100)
	for i := range arr {
		arr[i] = i
	}
	for i := 0; i < b.N; i++ {
		SliceN(arr, 10)
	}
}

func BenchmarkRandWeightSlice(b *testing.B) {
	weights := []int{10, 20, 30, 40}
	for i := 0; i < b.N; i++ {
		RandWeightSlice(weights)
	}
}

func BenchmarkRandWeightSliceTotal(b *testing.B) {
	weights := []int{10, 20, 30, 40}
	total := 100
	for i := 0; i < b.N; i++ {
		RandWeightSliceTotal(weights, total)
	}
}

func BenchmarkInTenThousandsProbability(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InTenThousandsProbability(5000)
	}
}

func BenchmarkInRandomProbability(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InRandomProbability(50, 100)
	}
}

func BenchmarkRangeInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RangeInts(0, 99, 10)
	}
}
