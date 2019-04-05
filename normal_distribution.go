package anomalia

type NormalDistribution struct {
	EpsilonThreshold float64
}

func (nd *NormalDistribution) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := nd.computeScores(timeSeries)
	return scoreList
}

func (nd *NormalDistribution) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	mean := timeSeries.Average()
	std := timeSeries.Stdev()

	var (
		score  float64
		scores []float64 = make([]float64, 0, len(timeSeries.Values))
	)

	for _, value := range timeSeries.Values {
		score = Pdf(mean, std)(value)
		if score >= nd.EpsilonThreshold {
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
