package wrand

import "math/rand"

const (
	choiceLen   = 2
	itemIndex   = 0
	weightIndex = 1
)

func NewRandChoices[T integer](choices [][]T) []RandChoice[T, T] {
	if len(choices) == 0 {
		return nil
	}

	ret := make([]RandChoice[T, T], 0, len(choices))
	for _, v := range choices {
		if len(v) < choiceLen {
			continue
		}
		ret = append(ret, NewRandChoice(v[itemIndex], v[weightIndex]))
	}
	return ret
}

func (c RandChooser[T, W]) PickN(n int) []T {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []T{c.Pick()}
	}

	var vals []T
	if n >= len(c.data) {
		vals := make([]T, 0, len(c.data))
		for _, v := range c.data {
			vals = append(vals, v.Item)
		}
		return vals
	}

	picked := make(map[int]struct{})
	maxWeight := c.max
	for range n {
		r := rand.Intn(maxWeight) + 1
		for idx := range c.totals {
			if _, ok := picked[idx]; ok {
				continue
			}

			weight := int(c.data[idx].Weight)
			if r > weight {
				r -= weight
				continue
			}

			picked[idx] = struct{}{}
			vals = append(vals, c.data[idx].Item)
			maxWeight -= weight
			break
		}
	}
	return vals
}
