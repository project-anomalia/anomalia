package anomalia

// AbsoluteThreshold holds absolute threshold algorithm configuration.
// It takes the difference of lower and upper thresholds with the current value as anomaly score.
type AbsoluteThreshold struct {
	LowerThreshold float64
	UpperThreshold float64
}

// NewAbsoluteThreshold returns AbsoluteThAbsoluteThreshold instance
func NewAbsoluteThreshold(lower, upper float64) *AbsoluteThreshold {
	return &AbsoluteThreshold{lower, upper}
}

// Run runs the absolute threshold algorithm over the time series
func (at *AbsoluteThreshold) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := at.computeScores(timeSeries)
	return scoreList
}

func (at *AbsoluteThreshold) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	scores := mapSlice(timeSeries.Values, func(value float64) float64 {
		if value > at.UpperThreshold {
			return value - at.UpperThreshold
		} else if value < at.LowerThreshold {
			return at.LowerThreshold - value
		} else {
			return 0.0
		}
	})
	scoreList := &ScoreList{timeSeries.Timestamps, scores}
	return scoreList, nil
}
