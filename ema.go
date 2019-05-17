package anomalia

import "math"

// ExponentialMovingAverage holds the algorithm configuration.
// It uses the value's deviation from the exponential moving average
// of a lagging window to determine anomalies scores.
type ExponentialMovingAverage struct {
	lagWindowSize   int
	smoothingFactor float64
}

// NewEma returns ExponentialMovingAverage instance
func NewEma() *ExponentialMovingAverage {
	return &ExponentialMovingAverage{2, 0.2}
}

// WithLagWindowSize sets the lagging window size.
func (ema *ExponentialMovingAverage) WithLagWindowSize(size int) *ExponentialMovingAverage {
	ema.lagWindowSize = size
	return ema
}

// WithSmoothingFactor sets the smoothing factor.
func (ema *ExponentialMovingAverage) WithSmoothingFactor(factor float64) *ExponentialMovingAverage {
	ema.smoothingFactor = factor
	return ema
}

// Run runs the exponential moving average algorithm over the time series
func (ema *ExponentialMovingAverage) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := ema.computeScores(timeSeries)
	return scoreList
}

func (ema *ExponentialMovingAverage) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	stdev := timeSeries.Stdev()
	scores := mapSliceWithIndex(timeSeries.Values, func(idx int, value float64) float64 {
		score := 0.0
		if idx < ema.lagWindowSize {
			score = computeScoresInLagWindow(timeSeries.Values[:idx+1], value, ema.smoothingFactor)
		} else {
			score = computeScoresInLagWindow(timeSeries.Values[idx-ema.lagWindowSize:idx+1], value, ema.smoothingFactor)
		}

		if stdev > 0.0 {
			score = score / stdev
		}
		return score
	})

	scoreList := &ScoreList{timeSeries.Timestamps, scores}
	return scoreList, nil
}

func computeScoresInLagWindow(data []float64, value, smoothingFactor float64) float64 {
	ema := Ema(data, smoothingFactor)[len(data)-1]
	return math.Abs(value - ema)
}
