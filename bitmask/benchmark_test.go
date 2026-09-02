package bitmask

import (
	"testing"
)

func BenchmarkSetBit_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetBit(0, 5)
	}
}

func BenchmarkSetBit_Int64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetBit(int64(0), 5)
	}
}

func BenchmarkResetBit_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ResetBit(0xFF, 5)
	}
}

func BenchmarkHasBit_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HasBit(0xFF, 5)
	}
}

func BenchmarkSetBitSlice(b *testing.B) {
	data := make([]byte, 16)
	for i := 0; i < b.N; i++ {
		SetBitSlice(data, 64)
	}
}

func BenchmarkResetBitSlice(b *testing.B) {
	data := make([]byte, 16)
	for i := 0; i < b.N; i++ {
		ResetBitSlice(data, 64)
	}
}

func BenchmarkHasBitSlice(b *testing.B) {
	data := make([]byte, 16)
	for i := 0; i < b.N; i++ {
		HasBitSlice(data, 64)
	}
}
