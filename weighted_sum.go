package anomalia

import "math"

// WeightedSum holds the weighted sum algorithm configuration.
// The weighted sum algorithm uses a weighted sum to calculate anomalies scores.
// It should be used ONLY on small datasets.
type WeightedSum struct {
	scoreWeight float64
	minEmaScore float64
	*ExponentialMovingAverage
	*Derivative
}

// NewWeightedSum returns weighted sum instance
func NewWeightedSum() *WeightedSum {
	return &WeightedSum{
		scoreWeight:              0.65,
		minEmaScore:              0.94,
		ExponentialMovingAverage: &ExponentialMovingAverage{2, 0.2},
		Derivative:               &Derivative{0.2},
	}
}

// WithScoreWeight sets Ema's score weight.
func (ws *WeightedSum) WithScoreWeight(weight float64) *WeightedSum {
	ws.scoreWeight = weight
	return ws
}

// WithMinEmaScore sets the minimal Ema score above which the weighted score is used.
func (ws *WeightedSum) WithMinEmaScore(value float64) *WeightedSum {
	ws.minEmaScore = value
	return ws
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
		weightedScore := emaScores[timestamp]*ws.scoreWeight + derivativeScores[timestamp]*(1-ws.scoreWeight)
		score := math.Max(emaScores[timestamp], weightedScore)

		if emaScores[timestamp] > ws.minEmaScore {
			return math.Max(score, derivativeScores[timestamp])
		}
		return score
	})

	scoreList := (&ScoreList{timeSeries.Timestamps, scores}).Denoise()
	return scoreList, nil
}
