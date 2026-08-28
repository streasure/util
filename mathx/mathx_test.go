package mathx

import (
	"testing"
)

func TestV2(t *testing.T) {
	v := NewV2(3, 4)
	if v.Len() != 5 {
		t.Errorf("V2.Len() = %f, want 5", v.Len())
	}
	if v.LenSqrt() != 25 {
		t.Errorf("V2.LenSqrt() = %f, want 25", v.LenSqrt())
	}
}

func TestV2Normalize(t *testing.T) {
	v := NewV2(3, 4)
	v.Normalize()
	if v.Len()-1.0 > 1e-10 {
		t.Errorf("V2.Normalize() length = %f, want 1.0", v.Len())
	}
}

func TestV2AddSub(t *testing.T) {
	v1 := NewV2(1, 2)
	v2 := NewV2(3, 4)
	result := Add(v1, v2)
	if result.X != 4 || result.Y != 6 {
		t.Errorf("Add = %v, want (4, 6)", result)
	}
	result = Sub(v2, v1)
	if result.X != 2 || result.Y != 2 {
		t.Errorf("Sub = %v, want (2, 2)", result)
	}
}

func TestIsFloatSame(t *testing.T) {
	if !IsFloatSame(1.0, 1.0) {
		t.Error("IsFloatSame(1.0, 1.0) should be true")
	}
	if IsFloatSame(1.0, 2.0) {
		t.Error("IsFloatSame(1.0, 2.0) should be false")
	}
}

func TestClamp(t *testing.T) {
	if Clamp(5.0, 1.0, 10.0) != 5.0 {
		t.Error("Clamp failed")
	}
	if Clamp(0.0, 1.0, 10.0) != 1.0 {
		t.Error("Clamp below min failed")
	}
	if Clamp(15.0, 1.0, 10.0) != 10.0 {
		t.Error("Clamp above max failed")
	}
}
