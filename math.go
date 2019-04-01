package anomalia

import "math"

// RoundFloat rounds float to closest int
func RoundFloat(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// Float64WithPrecision rounds float to certain precision
func Float64WithPrecision(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(RoundFloat(num*output)) / output
}
