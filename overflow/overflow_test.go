package overflow

import (
	"testing"
)

func TestCalcAddOverflow(t *testing.T) {
	addCount, overCount := CalcAddOverflow(100, 80, 30)
	if addCount != 20 || overCount != 10 {
		t.Errorf("CalcAddOverflow = (%d, %d), want (20, 10)", addCount, overCount)
	}

	addCount, overCount = CalcAddOverflow(100, 100, 10)
	if addCount != 0 || overCount != 10 {
		t.Errorf("CalcAddOverflow at max = (%d, %d), want (0, 10)", addCount, overCount)
	}
}

func TestCalcMinusOverflow(t *testing.T) {
	reduceCount, remainCount := CalcMinusOverflow(0, 80, 30)
	if reduceCount != 30 || remainCount != 0 {
		t.Errorf("CalcMinusOverflow = (%d, %d), want (30, 0)", reduceCount, remainCount)
	}

	reduceCount, remainCount = CalcMinusOverflow(0, 10, 30)
	if reduceCount != 10 || remainCount != 20 {
		t.Errorf("CalcMinusOverflow below min = (%d, %d), want (10, 20)", reduceCount, remainCount)
	}
}
