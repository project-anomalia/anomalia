package anomalia

import "math"

// PearsonCorrelator struct which holds the current and target time series.
type PearsonCorrelator struct {
	current, target *TimeSeries
}

// NewPearsonCorrelator returns an instance of the pearson correlator.
// It measures the linear correlation between the current and target time series.
// It should be used when the two time series are normally distributed.
//
// The correlation coefficient always has a value between -1 and +1 where:
//  - +1 is total positive linear correlation
//  - 0 is no linear correlation
//  - −1 is total negative linear correlation
//
// For the used formula, check: https://en.wikipedia.org/wiki/Pearson_correlation_coefficient
func NewPearsonCorrelator(current, target *TimeSeries) *PearsonCorrelator {
	return &PearsonCorrelator{current, target}
}

// Run runs the pearson correlation on the current and target time series.
// It returns the correlation coefficient which always has a value between -1 and +1.
func (pc *PearsonCorrelator) Run() float64 {
	pc.sanityCheck()

	currentSquares, targetSquares := sumOfSquares(pc.current.Values), sumOfSquares(pc.target.Values)
	currentAvg, targetAvg := Average(pc.current.Values), Average(pc.target.Values)
	n := float64(pc.current.Size())
	denom := math.Sqrt((currentSquares - n*currentAvg*currentAvg) * (targetSquares - n*targetAvg*targetAvg))

	if denom == 0 {
		return denom
	}
	return (sumOfProducts(pc.current.Values, pc.target.Values) - n*currentAvg*targetAvg) / denom
}

func (pc *PearsonCorrelator) sanityCheck() {
	if pc.current.Size() != pc.target.Size() {
		panic("current and target series do not have the same dimension")
	}
}
