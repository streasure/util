package mathutil

import (
	"cmp"
	"math"
	"math/rand"
)

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func Clamp[T Ordered](value, min, max T) T {
	if min > max {
		min, max = max, min
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func Max[T cmp.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Min[T cmp.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Abs[T Signed](value T) T {
	if value < 0 {
		return -value
	}
	return value
}

func Operator[T cmp.Ordered](operator string, v1, v2 T) bool {
	switch operator {
	case ">":
		return v1 > v2
	case "<":
		return v1 < v2
	case ">=":
		return v1 >= v2
	case "<=":
		return v1 <= v2
	case "=":
		return v1 == v2
	case "!=":
		return v1 != v2
	default:
		return false
	}
}

func RandInt(min, max int) int {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return min + rand.Intn(max-min+1)
}

func RandInt32(min, max int32) int32 {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return min + rand.Int31n(max-min+1)
}

func RandInt64(min, max int64) int64 {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return min + rand.Int63n(max-min+1)
}

// Deprecated: Use Clamp[T] instead.
func ClampInt(value, min, max int) int { return Clamp(value, min, max) }

// Deprecated: Use Clamp[T] instead.
func ClampInt32(value, min, max int32) int32 { return Clamp(value, min, max) }

// Deprecated: Use Clamp[T] instead.
func ClampInt64(value, min, max int64) int64 { return Clamp(value, min, max) }

// Deprecated: Use Clamp[T] instead.
func ClampFloat64(value, min, max float64) float64 { return Clamp(value, min, max) }

// Deprecated: Use Max[T] instead.
func MaxInt(x, y int) int { return Max(x, y) }

// Deprecated: Use Max[T] instead.
func MaxInt32(x, y int32) int32 { return Max(x, y) }

// Deprecated: Use Max[T] instead.
func MaxInt64(x, y int64) int64 { return Max(x, y) }

// Deprecated: Use Min[T] instead.
func MinInt(x, y int) int { return Min(x, y) }

// Deprecated: Use Min[T] instead.
func MinInt32(x, y int32) int32 { return Min(x, y) }

// Deprecated: Use Min[T] instead.
func MinInt64(x, y int64) int64 { return Min(x, y) }

// Deprecated: Use Abs[T] instead.
func AbsInt(value int) int { return Abs(value) }

// Deprecated: Use Abs[T] instead.
func AbsInt32(value int32) int32 { return Abs(value) }

// Deprecated: Use Abs[T] instead.
func AbsInt64(value int64) int64 { return Abs(value) }

// Deprecated: Use Operator[T] instead.
func OperatorInt(operator string, v1, v2 int) bool { return Operator(operator, v1, v2) }

// Deprecated: Use Operator[T] instead.
func OperatorInt32(operator string, v1, v2 int32) bool { return Operator(operator, v1, v2) }

// Deprecated: Use Operator[T] instead.
func OperatorInt64(operator string, v1, v2 int64) bool { return Operator(operator, v1, v2) }

func ClampFloat64Math(value, min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	return math.Min(math.Max(value, min), max)
}
