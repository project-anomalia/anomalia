package anomalia

import "math"

// Average returns the average of the input
func Average(input []float64) float64 {
	return Sum(input) / float64(len(input))
}

// Sum returns the sum of all elements in the input
func Sum(input []float64) float64 {
	var sum float64
	for _, value := range input {
		sum += value
	}
	return sum
}

// Variance returns the variance of the input
func Variance(input []float64) (variance float64) {
	avg := Average(input)
	for _, value := range input {
		variance += (value - avg) * (value - avg)
	}
	return variance / float64(len(input))
}

// Stdev returns the standard deviation of the input
func Stdev(input []float64) float64 {
	variance := Variance(input)
	return math.Pow(variance, 0.5)
}

// RoundFloat rounds float to closest int
func RoundFloat(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// Float64WithPrecision rounds float to certain precision
func Float64WithPrecision(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(RoundFloat(num*output)) / output
}

// Pdf returns the probability density function
func Pdf(mean, stdev float64) func(float64) float64 {
	return func(x float64) float64 {
		numexp := math.Pow(x-mean, 2) / (2 * math.Pow(stdev, 2))
		denom := stdev * math.Sqrt(2*math.Pi)
		numer := math.Pow(math.E, numexp*-1)
		return numer / denom
	}
}
