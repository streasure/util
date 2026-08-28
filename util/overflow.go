package util

const (
	Percentage = 10000
)

func CalcAddOverflow(max, cur, add int64) (addCount int64, overCount int64) {
	if add <= 0 {
		return 0, 0
	}
	if cur >= max {
		return 0, add
	}

	switch {
	case cur < 0:
		newCount := add + cur
		if newCount > max {
			return max - cur, newCount - max
		}
		return add, 0
	default:
		newCount := uint64(add) + uint64(cur)
		if newCount > uint64(max) {
			return max - cur, int64(newCount - uint64(max))
		}
		return add, 0
	}
}

func CalcMinusOverflow(min, cur, minus int64) (reduceCount int64, remainCount int64) {
	if minus < 0 {
		return 0, 0
	}
	if cur <= min {
		return 0, minus
	}

	switch {
	case cur < 0:
		tmpCount := uint64(-cur) + uint64(minus)
		if tmpCount > uint64(-min) {
			return cur - min, int64(tmpCount - uint64(-min))
		}
		return minus, 0
	default:
		newCount := cur - minus
		if newCount < min {
			return cur - min, min - newCount
		}
		return minus, 0
	}
}
