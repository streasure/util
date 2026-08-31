package util

import (
	"cmp"
	"sort"
)

type RangePair[T cmp.Ordered] struct {
	Start T
	End   T
}

func CheckRangeIntersect[T cmp.Ordered](ranges []RangePair[T]) bool {
	if len(ranges) < 2 {
		return false
	}
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start < ranges[j].Start
	})

	n := len(ranges)
	for i := 1; i < n; i++ {
		if ranges[i].Start < ranges[i-1].End {
			return true
		}
	}

	return false
}
