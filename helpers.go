package anomalia

import "sort"

type Mapper func(float64) float64
type Predicate func(float64) bool

func MinMax(data []float64) (float64, float64) {
	var (
		max float64 = data[0]
		min float64 = data[0]
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

func Map(slice []float64, mapper Mapper) []float64 {
	for idx, value := range slice {
		slice[idx] = mapper(value)
	}
	return slice
}

func Filter(slice []float64, predicate Predicate) (ret []float64) {
	for _, value := range slice {
		if predicate(value) {
			ret = append(ret, value)
		}
	}
	return
}

func CopySlice(input []float64) []float64 {
	s := make([]float64, len(input))
	copy(s, input)
	return s
}

func SortedCopy(input []float64) (copy []float64) {
	copy = CopySlice(input)
	sort.Float64s(copy)
	return
}
