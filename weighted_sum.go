package anomalia

import "math"

// WeightedSum holds the weighted sum algorithm configuration.
// The weighted sum algorithm uses a weighted sum to calculate anomalies scores.
// It should be used ONLY on small datasets
type WeightedSum struct {
	EmaWeight      float64
	EmaSignificant float64
	*ExponentialMovingAverage
	*Derivative
}

// NewWeightedSum returns weighted sum instance
func NewWeightedSum() *WeightedSum {
	return &WeightedSum{
		EmaWeight:                0.65,
		EmaSignificant:           0.94,
		ExponentialMovingAverage: &ExponentialMovingAverage{2, 0.2},
		Derivative:               &Derivative{0.2},
	}
}

// Run runs the weighted sum algorithm over the time series
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
