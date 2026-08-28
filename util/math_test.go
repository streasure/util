package util

import (
	"testing"
)

func TestClamp(t *testing.T) {
	// Test with int
	if Clamp(5, 1, 10) != 5 {
		t.Error("Clamp int failed")
	}
	if Clamp(0, 1, 10) != 1 {
		t.Error("Clamp int below min failed")
	}
	if Clamp(15, 1, 10) != 10 {
		t.Error("Clamp int above max failed")
	}

	// Test with float64
	if Clamp(5.5, 1.0, 10.0) != 5.5 {
		t.Error("Clamp float64 failed")
	}
	if Clamp(0.0, 1.0, 10.0) != 1.0 {
		t.Error("Clamp float64 below min failed")
	}

	// Test with reversed min/max
	if Clamp(5, 10, 1) != 5 {
		t.Error("Clamp reversed min/max failed")
	}
}

func TestMax(t *testing.T) {
	if Max(3, 5) != 5 {
		t.Error("Max int failed")
	}
	if Max(5.5, 3.3) != 5.5 {
		t.Error("Max float64 failed")
	}
	if Max("abc", "xyz") != "xyz" {
		t.Error("Max string failed")
	}
}

func TestMin(t *testing.T) {
	if Min(3, 5) != 3 {
		t.Error("Min int failed")
	}
	if Min(5.5, 3.3) != 3.3 {
		t.Error("Min float64 failed")
	}
}

func TestAbs(t *testing.T) {
	if Abs(-5) != 5 {
		t.Error("Abs int failed")
	}
	if Abs(int32(-5)) != 5 {
		t.Error("Abs int32 failed")
	}
	if Abs(int64(-5)) != 5 {
		t.Error("Abs int64 failed")
	}
}

func TestOperator(t *testing.T) {
	if !Operator(">", 5, 3) {
		t.Error("Operator > failed")
	}
	if !Operator("=", 5, 5) {
		t.Error("Operator = failed")
	}
	if !Operator("!=", 5, 3) {
		t.Error("Operator != failed")
	}
	if !Operator("<", 3, 5) {
		t.Error("Operator < failed")
	}
	if !Operator(">=", 5, 5) {
		t.Error("Operator >= failed")
	}
	if !Operator("<=", 3, 5) {
		t.Error("Operator <= failed")
	}
}

func TestRandInt(t *testing.T) {
	for i := 0; i < 1000; i++ {
		got := RandInt(5, 10)
		if got < 5 || got > 10 {
			t.Errorf("RandInt(5, 10) = %d, want [5, 10]", got)
		}
	}
}

func TestDeprecatedWrappers(t *testing.T) {
	// Test backward compatibility
	if ClampInt(5, 1, 10) != 5 {
		t.Error("ClampInt failed")
	}
	if ClampFloat64(5.5, 1.0, 10.0) != 5.5 {
		t.Error("ClampFloat64 failed")
	}
	if MaxInt(3, 5) != 5 {
		t.Error("MaxInt failed")
	}
	if MinInt(3, 5) != 3 {
		t.Error("MinInt failed")
	}
	if AbsInt(-5) != 5 {
		t.Error("AbsInt failed")
	}
	if !OperatorInt(">", 5, 3) {
		t.Error("OperatorInt failed")
	}
}
