package util

import (
	"testing"
)

func TestSetBit(t *testing.T) {
	// Test with int32
	var data32 int32
	data32 = SetBit(data32, 0)
	if !HasBit(data32, 0) {
		t.Error("SetBit int32(0) failed")
	}
	data32 = SetBit(data32, 31)
	if !HasBit(data32, 31) {
		t.Error("SetBit int32(31) failed")
	}

	// Test with int64
	var data64 int64
	data64 = SetBit(data64, 0)
	if !HasBit(data64, 0) {
		t.Error("SetBit int64(0) failed")
	}
	data64 = SetBit(data64, 63)
	if !HasBit(data64, 63) {
		t.Error("SetBit int64(63) failed")
	}

	// Test with uint32
	var dataUint32 uint32
	dataUint32 = SetBit(dataUint32, 5)
	if !HasBit(dataUint32, 5) {
		t.Error("SetBit uint32(5) failed")
	}
}

func TestResetBit(t *testing.T) {
	var data int32
	data = SetBit(data, 5)
	data = ResetBit(data, 5)
	if HasBit(data, 5) {
		t.Error("ResetBit failed")
	}
}

func TestHasBit(t *testing.T) {
	var data int32
	if HasBit(data, 0) {
		t.Error("HasBit should be false for zero value")
	}
	data = SetBit(data, 10)
	if !HasBit(data, 10) {
		t.Error("HasBit should be true after SetBit")
	}
}

func TestBitSlice(t *testing.T) {
	data := make([]byte, 2)
	data = SetBitSlice(data, 0)
	if !HasBitSlice(data, 0) {
		t.Error("SetBitSlice(0) failed")
	}
	data = SetBitSlice(data, 9)
	if !HasBitSlice(data, 9) {
		t.Error("SetBitSlice(9) failed")
	}

	data = ResetBitSlice(data, 0)
	if HasBitSlice(data, 0) {
		t.Error("ResetBitSlice(0) failed")
	}
}

func TestDeprecatedBitWrappers(t *testing.T) {
	// Test backward compatibility
	var data32 int32
	data32 = SetBit32(data32, 0)
	if !HasBit32(data32, 0) {
		t.Error("SetBit32 failed")
	}

	var data64 int64
	data64 = SetBit64(data64, 0)
	if !HasBit64(data64, 0) {
		t.Error("SetBit64 failed")
	}
}
