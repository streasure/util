package overflow

import (
	"testing"
)

func BenchmarkCalcAddOverflow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalcAddOverflow(100, 50, 30)
	}
}

func BenchmarkCalcAddOverflow_Overflow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalcAddOverflow(100, 80, 50)
	}
}

func BenchmarkCalcMinusOverflow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalcMinusOverflow(0, 50, 30)
	}
}

func BenchmarkCalcMinusOverflow_Overflow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalcMinusOverflow(0, 50, 80)
	}
}
