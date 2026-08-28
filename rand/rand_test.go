package rand

import (
	"testing"
)

func TestInTenThousandsProbability(t *testing.T) {
	if !InTenThousandsProbability(10000) {
		t.Error("InTenThousandsProbability(10000) should be true")
	}
	if InTenThousandsProbability(0) {
		t.Error("InTenThousandsProbability(0) should be false")
	}
}

func TestInRandomProbability(t *testing.T) {
	if !InRandomProbability(1, 1) {
		t.Error("InRandomProbability(1,1) should be true")
	}
	if InRandomProbability(0, 1) {
		t.Error("InRandomProbability(0,1) should be false")
	}
}

func TestRangeInt(t *testing.T) {
	for i := 0; i < 1000; i++ {
		got := RangeInt(5, 10)
		if got < 5 || got > 10 {
			t.Errorf("RangeInt(5, 10) = %d, want [5, 10]", got)
		}
	}
}

func TestRangeInts(t *testing.T) {
	got := RangeInts(1, 10, 5)
	if len(got) != 5 {
		t.Errorf("RangeInts len = %d, want 5", len(got))
	}
	seen := make(map[int]bool)
	for _, v := range got {
		if seen[v] {
			t.Errorf("RangeInts returned duplicate %d", v)
		}
		seen[v] = true
	}
}

func TestSliceOne(t *testing.T) {
	arr := []int{1, 2, 3}
	val, ok := SliceOne(arr)
	if !ok || val < 1 || val > 3 {
		t.Errorf("SliceOne = %d, %v", val, ok)
	}

	_, ok = SliceOne([]int{})
	if ok {
		t.Error("SliceOne on empty should return false")
	}
}

func TestRandWeightSlice(t *testing.T) {
	weights := []int{10, 20, 70}
	idx := RandWeightSlice(weights)
	if idx < 0 || idx >= len(weights) {
		t.Errorf("RandWeightSlice = %d, want [0, %d]", idx, len(weights)-1)
	}
}

func TestWeightRandom(t *testing.T) {
	wr := NewWeightRandom[int]()
	wr.Add(10, 1)
	wr.Add(20, 2)
	wr.Add(70, 3)

	val, ok := wr.Rand()
	if !ok || val < 1 || val > 3 {
		t.Errorf("WeightRandom.Rand = %d, %v", val, ok)
	}

	vals := wr.RandN(2)
	if len(vals) != 2 {
		t.Errorf("WeightRandom.RandN len = %d, want 2", len(vals))
	}
}
