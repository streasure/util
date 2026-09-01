package mathx

import (
	"math"

	"github.com/streasure/util/mathutil"
)

const Epsilon float64 = 1e-10
const MinFloat = float64(1.1754943508222875e-38)

func IsFloatSame(f1 float64, f2 float64) bool {
	if f1 == f2 {
		return true
	}

	diff := math.Abs(f1 - f2)
	if f1 == 0 || f2 == 0 || diff < MinFloat {
		return diff < Epsilon*Epsilon
	}

	return diff/(math.Abs(f1)+math.Abs(f2)) < Epsilon
}

// Deprecated: Use util.Clamp[T] instead.
func Clamp(f float64, low float64, high float64) float64 {
	return mathutil.Clamp(f, low, high)
}

func IsClamped(f float64, low float64, high float64) bool {
	return f >= low && f <= high
}
