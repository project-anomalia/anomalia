package anomalia

type AbsoluteThreshold struct {
	LowerThreshold float64
	UpperThreshold float64
}

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
	scoreList := &ScoreList{
		Timestamps: timeSeries.Timestamps,
		Scores:     scores,
	}
	return scoreList, nil
}
