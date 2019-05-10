package anomalia

import (
	"sort"
	"sync"
)

type mapper func(float64) float64
type mapperWithIndex func(int, float64) float64
type predicate func(float64) bool

func minMax(data []float64) (float64, float64) {
	var (
		max = data[0]
		min = data[0]
	)
	for _, value := range data {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func mapSlice(slice []float64, m mapper) []float64 {
	var (
		wg     sync.WaitGroup
		result = make([]float64, len(slice))
	)

	wg.Add(len(slice))
	for i, value := range slice {
		go func(i int, value float64) {
			defer wg.Done()
			result[i] = m(value)
		}(i, value)
	}
	wg.Wait()

	return result
}

func mapSliceWithIndex(slice []float64, m mapperWithIndex) []float64 {
	var (
		wg     sync.WaitGroup
		result = make([]float64, len(slice))
	)

	wg.Add(len(slice))
	for idx, value := range slice {
		go func(idx int, value float64) {
			defer wg.Done()
			result[idx] = m(idx, value)
		}(idx, value)
	}
	wg.Wait()

	return result
}

func filter(slice []float64, predicate predicate) (ret []float64) {
	for _, value := range slice {
		if predicate(value) {
			ret = append(ret, value)
		}
	}
	return
}

func copySlice(input []float64) []float64 {
	s := make([]float64, len(input))
	copy(s, input)
	return s
}

func sortedCopy(input []float64) (copy []float64) {
	copy = copySlice(input)
	sort.Float64s(copy)
	return
}

func insertAt(slice []float64, pos int, elem float64) []float64 {
	if pos < 0 {
		pos = 0
	} else if pos >= len(slice) {
		pos = len(slice)
	}
	out := make([]float64, len(slice)+1)
	copy(out[:pos], slice[:pos])
	out[pos] = elem
	copy(out[pos+1:], slice[pos:])
	return out
}

func mapIntKeys(dict map[int]float64) []int {
	keys := make([]int, len(dict))
	i := 0
	for key := range dict {
		keys[i] = key
		i++
	}
	sort.Ints(keys)
	return keys
}

func mapFloat64Keys(m map[float64]float64) []float64 {
	keys := make([]float64, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	sort.Float64s(keys)
	return keys
}

func indexOf(slice []float64, value float64) int {
	for idx := range slice {
		if slice[idx] == value {
			return idx
		}
	}
	return -1
}

func unpackMap(m map[float64]float64) ([]float64, []float64) {
	keys := mapFloat64Keys(m)
	values := make([]float64, len(keys))
	for idx, key := range keys {
		values[idx] = m[key]
	}
	return keys, values
}

func sumOfSquares(s []float64) float64 {
	sum := 0.0
	for _, val := range s {
		sum += val * val
	}
	return sum
}

func sumOfProducts(s1 []float64, s2 []float64) float64 {
	sum := 0.0
	for i := range s1 {
		sum += s1[i] * s2[i]
	}
	return sum
}
