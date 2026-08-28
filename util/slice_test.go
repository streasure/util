package util

import (
	"testing"
)

func TestStrToSlice(t *testing.T) {
	str := "hello"
	b := StrToSlice(str)
	if string(b) != "hello" {
		t.Errorf("StrToSlice = %q, want %q", string(b), "hello")
	}
}

func TestUniqueSlice(t *testing.T) {
	input := []int{1, 2, 2, 3, 3, 3}
	got := UniqueSlice(input)
	if len(got) != 3 {
		t.Errorf("UniqueSlice len = %d, want 3", len(got))
	}
}

func TestEqualSlice(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{3, 2, 1}
	if !EqualSlice(a, b) {
		t.Error("EqualSlice should be true for same elements")
	}
	c := []int{1, 2, 4}
	if EqualSlice(a, c) {
		t.Error("EqualSlice should be false for different elements")
	}
}

func TestIsNil(t *testing.T) {
	if !IsNil(nil) {
		t.Error("IsNil(nil) should be true")
	}
	var p *int
	if !IsNil(p) {
		t.Error("IsNil(nil pointer) should be true")
	}
	v := 5
	if IsNil(&v) {
		t.Error("IsNil(non-nil pointer) should be false")
	}
}
