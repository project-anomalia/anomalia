package anomalia

// AbsoluteThreshold holds absolute threshold algorithm configuration.
// It takes the difference of lower and upper thresholds with the current value as anomaly score.
type AbsoluteThreshold struct {
	lowerThreshold float64
	upperThreshold float64
}

// NewAbsoluteThreshold returns AbsoluteThAbsoluteThreshold instance.
func NewAbsoluteThreshold() *AbsoluteThreshold {
	return &AbsoluteThreshold{}
}

// WithBounds sets both lower and upper thresholds.
func (at *AbsoluteThreshold) WithBounds(lower, upper float64) *AbsoluteThreshold {
	at.lowerThreshold = lower
	at.upperThreshold = upper
	return at
}

// Run runs the absolute threshold algorithm over the time series.
func (at *AbsoluteThreshold) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := at.computeScores(timeSeries)
	return scoreList
}

func (at *AbsoluteThreshold) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	scores := mapSlice(timeSeries.Values, func(value float64) float64 {
		if value > at.upperThreshold {
			return value - at.upperThreshold
		} else if value < at.lowerThreshold {
			return at.lowerThreshold - value
		} else {
			return 0.0
		}
	})
	scoreList := &ScoreList{timeSeries.Timestamps, scores}
	return scoreList, nil
}
