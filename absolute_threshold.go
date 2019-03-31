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
	var (
		score  float64
		scores []float64 = make([]float64, 0, len(timeSeries.Values))
	)
	for _, value := range timeSeries.Values {
		if value > at.UpperThreshold {
			score = value - at.UpperThreshold
		} else if value < at.LowerThreshold {
			score = at.LowerThreshold - value
		} else {
			score = 0.0
		}
		scores = append(scores, score)
	}

	scoreList := &ScoreList{
		Timestamps: timeSeries.Timestamps,
		Scores:     scores,
	}
	return scoreList, nil
}
