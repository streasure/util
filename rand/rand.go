package rand

import "math/rand/v2"

const TenThousand = 10000

func InTenThousandsProbability(rate int) bool {
	return InRandomProbability(rate, TenThousand)
}

func InRandomProbability(rate int, total int) bool {
	if rate <= 0 {
		return false
	}
	if rate >= total {
		return true
	}
	return IntN(total) < rate
}

func RangeInt32(min, max int32) int32 {
	if min >= max {
		return min
	}
	return min + Int32N(max-min+1)
}

func RangeInt64(min, max int64) int64 {
	if min >= max {
		return min
	}
	return min + rand.Int64N(max-min+1)
}

func RangeInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + IntN(max-min+1)
}

func RangeInts(min, max, n int) []int {
	if min > max || n <= 0 {
		return nil
	}

	ret := generateInts(min, max-min+1)
	if n > max-min+1 {
		return ret
	}

	shuffleN(n, ret)
	return ret[:n]
}

func Int32N(n int32) int32 {
	return rand.Int32N(n)
}

func IntN(n int) int {
	return rand.IntN(n)
}

func Float64() float64 {
	return rand.Float64()
}

func SliceOne[T any](array []T) (T, bool) {
	if len(array) == 0 {
		var zero T
		return zero, false
	}
	return array[IntN(len(array))], true
}

func SliceN[T any](array []T, n int) []T {
	if n <= 0 {
		return nil
	}
	if len(array) <= n {
		return array
	}

	indexArray := generateInts(0, len(array))
	shuffleN(n, indexArray)

	ret := make([]T, 0, n)
	for i := 0; i < n; i++ {
		ret = append(ret, array[indexArray[i]])
	}
	return ret
}

func RandWeightSlice(weights []int) int {
	if len(weights) == 0 {
		return -1
	}
	if len(weights) == 1 {
		return 0
	}

	var totalWeight int
	for _, weight := range weights {
		if weight < 0 {
			return -1
		}
		totalWeight += weight
	}

	if totalWeight <= 0 {
		return -1
	}

	return RandWeightSliceTotal(weights, totalWeight)
}

func RandWeightSliceTotal(weights []int, totalWeight int) int {
	if totalWeight <= 0 {
		return -1
	}

	random := IntN(totalWeight)
	var sum int
	for k, v := range weights {
		sum += v
		if random < sum {
			return k
		}
	}
	return -1
}

func generateInts(start, n int) []int {
	if n <= 0 {
		return nil
	}
	ret := make([]int, n)
	for i := 0; i < n; i++ {
		ret[i] = start + i
	}
	return ret
}

func shuffleN[T any](n int, array []T) {
	if len(array) <= n {
		return
	}
	rand.Shuffle(len(array), func(i, j int) {
		array[i], array[j] = array[j], array[i]
	})
}
