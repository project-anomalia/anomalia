package anomalia

import "math"

// CrossCorrelator holds cross correlator algorithm parameters and settings.
// It is calculated by multiplying and summing the current and target time series together.
//
// This implementation uses normalized time series which makes scoring easy to understand:
// 	- The higher the coefficient, the higher the correlation is.
// 	- The maximum value of the correlation coefficient is 1.
//	- The minimum value of the correlation coefficient is -1.
//	- Two time series are exactly the same when their correlation coefficient is equal to 1.
type CrossCorrelator struct {
	current, target *TimeSeries
	maxShift        float64
	impact          float64
}

// CorrelationResult holds detected correlation result.
type CorrelationResult struct {
	Shift              float64
	Coefficient        float64
	ShiftedCoefficient float64
}

// NewCrossCorrelator returns an instance of the cross correlator.
func NewCrossCorrelator(current *TimeSeries, target *TimeSeries) *CrossCorrelator {
	return &CrossCorrelator{
		current:  current,
		target:   target,
		maxShift: 60 * 1000,
		impact:   0.05,
	}
}

// WithMaxShift sets the maximal shift in seconds.
func (cc *CrossCorrelator) WithMaxShift(shift float64) *CrossCorrelator {
	cc.maxShift = shift * 1000
	return cc
}

// WithImpact sets impact of shift on shifted correlation coefficient.
func (cc *CrossCorrelator) WithImpact(impact float64) *CrossCorrelator {
	cc.impact = impact
	return cc
}

// GetCorrelationResult runs the cross correlation algorithm.
func (cc *CrossCorrelator) GetCorrelationResult() CorrelationResult {
	cc.sanityCheck()
	return cc.detectCorrelation()
}

// Run runs the cross correlation algorithm and returns only the coefficient.
func (cc *CrossCorrelator) Run() float64 {
	return cc.GetCorrelationResult().Coefficient
}

func (cc *CrossCorrelator) detectCorrelation() CorrelationResult {
	cc.current, cc.target = cc.current.Normalize(), cc.target.Normalize()
	cc.current.Align(cc.target)

	correlations := make([][]float64, 0)
	shiftedCorrelations := make([]float64, 0)

	currentValues, targetValues := cc.current.Values, cc.target.Values
	currentAvg, targetAvg := cc.current.Average(), cc.target.Average()
	currentStdev, targetStdev := cc.current.Stdev(), cc.target.Stdev()

	n := cc.current.Size()
	denom := currentStdev * targetStdev * float64(n)
	allowedShiftStep := findMaxAllowedShift(cc.current.Timestamps, cc.maxShift)

	var shiftLowerBound, shiftUpperBound int
	if allowedShiftStep != -1 {
		shiftLowerBound = -allowedShiftStep
		shiftUpperBound = allowedShiftStep
	} else {
		shiftLowerBound = 0
		shiftUpperBound = 1
	}

	for delay := shiftLowerBound; delay < shiftUpperBound; delay++ {
		_delay := math.Abs(cc.current.Timestamps[AbsInt(delay)] - cc.current.Timestamps[0])
		sum := 0.0
		for i := 0; i < n; i++ {
			j := i + delay
			if j < 0 || j >= n {
				continue
			} else {
				sum += (currentValues[i] - currentAvg) * (targetValues[j] - targetAvg)
			}
		}

		// Calculate correlation coefficient
		r := sum
		if denom != 0 {
			r = sum / denom
		}
		correlations = append(correlations, []float64{_delay, r})

		// Take into account the maximal shift
		if cc.maxShift > 0 {
			r *= 1 + _delay/float64(cc.maxShift)*cc.impact
		}
		shiftedCorrelations = append(shiftedCorrelations, r)
	}

	maxCorrelation := findMaxCorrelation(correlations)
	_, maxShiftedCorrelation := minMax(shiftedCorrelations)
	return CorrelationResult{
		Shift:              maxCorrelation[0],
		Coefficient:        maxCorrelation[1],
		ShiftedCoefficient: maxShiftedCorrelation,
	}
}

func findMaxAllowedShift(timestamps []float64, target float64) int {
	initialTimestamp := timestamps[0]
	residualTimestamps := mapSlice(timestamps, func(ts float64) float64 {
		return ts - initialTimestamp
	})
	// Find the first element in timestamps whose value is bigger than target.
	pos := -1
	lowerBound, upperBound := 0, len(residualTimestamps)
	for lowerBound < upperBound {
		pos = int(lowerBound + (upperBound-lowerBound)/2)
		if timestamps[pos] > target {
			upperBound = pos
		} else {
			lowerBound = pos + 1
		}
	}
	return pos
}

func findMaxCorrelation(data [][]float64) []float64 {
	max := data[0]
	for _, slice := range data {
		if slice[1] > max[1] {
			max = slice
		}
	}
	return max
}

func (cc *CrossCorrelator) sanityCheck() {
	if cc.current.Size() < 2 || cc.target.Size() < 2 {
		panic("not enough data points")
	}
}
