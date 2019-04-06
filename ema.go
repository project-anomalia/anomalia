package anomalia

import (
	"math"
)

type ExponentialMovingAverage struct {
	LagWindowSize   int
	SmoothingFactor float64
}

func (ema *ExponentialMovingAverage) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := ema.computeScores(timeSeries)
	return scoreList
}

func (ema *ExponentialMovingAverage) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	stdev := timeSeries.Stdev()
	scores := mapSliceWithIndex(timeSeries.Values, func(idx int, value float64) float64 {
		score := 0.0
		if idx < ema.LagWindowSize {
			score = computeScoresInLagWindow(timeSeries.Values[:idx+1], value, ema.SmoothingFactor)
		} else {
			score = computeScoresInLagWindow(timeSeries.Values[idx-ema.LagWindowSize:idx+1], value, ema.SmoothingFactor)
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
