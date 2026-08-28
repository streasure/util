package rand

import (
	"fmt"
	"strings"
)

type WeightRandom[T any] struct {
	totalWeight int32
	weights     []int32
	items       []T
}

func NewWeightRandom[T any]() *WeightRandom[T] {
	return &WeightRandom[T]{}
}

func (s *WeightRandom[T]) Add(weight int32, item T) {
	if weight > 0 {
		s.totalWeight += weight
		s.weights = append(s.weights, weight)
		s.items = append(s.items, item)
	}
}

func (s *WeightRandom[T]) Rand() (T, bool) {
	if s == nil || s.totalWeight <= 0 {
		var zero T
		return zero, false
	}

	r := RangeInt32(1, s.totalWeight)
	for index, weight := range s.weights {
		if r <= weight {
			return s.items[index], true
		}
		r -= weight
	}
	var ret T
	return ret, false
}

func (s *WeightRandom[T]) RandN(n int) []T {
	if s == nil || s.totalWeight <= 0 || n <= 0 {
		return nil
	}

	ret := make([]T, 0, n)
	if len(s.items) <= n {
		ret = append(ret, s.items...)
		return ret
	}

	picked := make(map[int]struct{})
	totalWeight := s.totalWeight

	for i := 0; i < n; i++ {
		r := RangeInt32(1, totalWeight)
		for index, weight := range s.weights {
			if _, ok := picked[index]; ok {
				continue
			}
			if r > weight {
				r -= weight
				continue
			}
			picked[index] = struct{}{}
			ret = append(ret, s.items[index])
			totalWeight -= weight
			break
		}
	}
	return ret
}

func (s *WeightRandom[T]) RandOneExclude(fn func(t T) bool) (T, bool) {
	if s == nil || s.totalWeight <= 0 {
		var zero T
		return zero, false
	}
	picked := make(map[int]struct{})
	totalWeight := s.totalWeight

	for i := 0; i < len(s.items); i++ {
		if totalWeight <= 0 {
			break
		}
		r := RangeInt32(1, totalWeight)
		for index, weight := range s.weights {
			if _, ok := picked[index]; ok {
				continue
			}
			if r > weight {
				r -= weight
				continue
			}
			picked[index] = struct{}{}
			totalWeight -= weight
			if fn(s.items[index]) {
				return s.items[index], true
			}
		}
	}
	var zero T
	return zero, false
}

func (s *WeightRandom[T]) RandNExclude(n int, fn func(t T) bool) []T {
	if s == nil || s.totalWeight <= 0 || n <= 0 {
		return nil
	}
	ret := make([]T, 0, n)
	picked := make(map[int]struct{})
	totalWeight := s.totalWeight

	var pickedCount int
	for i := 0; i < len(s.items); i++ {
		if pickedCount >= n {
			break
		}
		if totalWeight <= 0 {
			break
		}
		r := RangeInt32(1, totalWeight)
		for index, weight := range s.weights {
			if _, ok := picked[index]; ok {
				continue
			}
			if r > weight {
				r -= weight
				continue
			}
			picked[index] = struct{}{}
			totalWeight -= weight
			if fn(s.items[index]) {
				ret = append(ret, s.items[index])
				pickedCount++
			}
			break
		}
	}
	return ret
}

func (s *WeightRandom[T]) String() string {
	if s == nil {
		return "nil"
	}
	strs := make([]string, 0, len(s.items))
	for i := 0; i < len(s.weights); i++ {
		strs = append(strs, fmt.Sprintf("%v:%v", s.items[i], s.weights[i]))
	}
	return strings.Join(strs, " ")
}
