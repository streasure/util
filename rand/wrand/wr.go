package wrand

import (
	"errors"
	"math/rand"
	"sort"
	"strconv"
)

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

const (
	maxInt = 1<<(strconv.IntSize-1) - 1
)

var (
	errWeightOverflow = errors.New("sum of RandChoice Weights exceeds max int")
	errNoValidChoices = errors.New("zero Choices with Weight >= 1")
)

type RandChoice[T any, W integer] struct {
	Item   T
	Weight W
}

func NewRandChoice[T any, W integer](item T, weight W) RandChoice[T, W] {
	return RandChoice[T, W]{Item: item, Weight: weight}
}

type RandChooser[T any, W integer] struct {
	data   []RandChoice[T, W]
	totals []int
	max    int
}

func NewRandChooser[T any, W integer](choices ...RandChoice[T, W]) (*RandChooser[T, W], error) {
	sort.Slice(choices, func(i, j int) bool {
		return choices[i].Weight < choices[j].Weight
	})

	totals := make([]int, len(choices))
	runningTotal := 0
	for i, c := range choices {
		weight := int(c.Weight)
		if weight < 0 {
			continue
		}

		if (maxInt-runningTotal) <= weight {
			return nil, errWeightOverflow
		}
		runningTotal += weight
		totals[i] = runningTotal
	}

	if runningTotal < 1 {
		return nil, errNoValidChoices
	}

	return &RandChooser[T, W]{data: choices, totals: totals, max: runningTotal}, nil
}

func (c RandChooser[T, W]) Pick() T {
	r := rand.Intn(c.max) + 1
	i := searchInts(c.totals, r)
	return c.data[i].Item
}

func (c RandChooser[T, W]) PickSource(rs *rand.Rand) T {
	r := rs.Intn(c.max) + 1
	i := searchInts(c.totals, r)
	return c.data[i].Item
}

func searchInts(a []int, x int) int {
	i, j := 0, len(a)
	for i < j {
		h := int(uint(i+j) >> 1)
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
