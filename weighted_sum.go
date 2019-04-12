package anomalia

import "math"

type WeightedSum struct {
	EmaWeight      float64
	EmaSignificant float64
	*ExponentialMovingAverage
	*Derivative
}

func (ws *WeightedSum) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := ws.computeScores(timeSeries)
	return scoreList
}

func (ws *WeightedSum) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	emaScores := ws.ExponentialMovingAverage.Run(timeSeries).Zip()
	derivativeScores := ws.Derivative.Run(timeSeries).Zip()

	scores := mapSlice(timeSeries.Timestamps, func(timestamp float64) float64 {
		weightedScore := emaScores[timestamp]*ws.EmaWeight + derivativeScores[timestamp]*(1-ws.EmaWeight)
		score := math.Max(emaScores[timestamp], weightedScore)

		if emaScores[timestamp] > ws.EmaSignificant {
			return math.Max(score, derivativeScores[timestamp])
		}
		return score
	})

	scoreList := (&ScoreList{timeSeries.Timestamps, scores}).Denoise()
	return scoreList, nil
}
