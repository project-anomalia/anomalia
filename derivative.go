package anomalia

import "math"

// Derivative holds the derivative algorithm configuration.
// It uses the derivative of the current value as anomaly score.
type Derivative struct {
	SmoothingFactor float64
}

// NewDerivative return Derivative instance
func NewDerivative(smoothingFactor float64) *Derivative {
	return &Derivative{smoothingFactor}
}

// Run runs the derivative algorithm over the time series
func (d *Derivative) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := d.computeScores(timeSeries)
	return scoreList
}

func (d *Derivative) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	derivatives := d.computeDerivatives(timeSeries)
	derivativesEma := Ema(derivatives, d.SmoothingFactor)

	scores := mapSliceWithIndex(timeSeries.Values, func(i int, value float64) float64 {
		return math.Abs(derivatives[i] - derivativesEma[i])
	})

	stdev := Stdev(scores)
	if stdev != 0.0 {
		scores = mapSlice(scores, func(score float64) float64 {
			return score / stdev
		})
	}
	scoreList := (&ScoreList{timeSeries.Timestamps, scores}).Denoise()
	return scoreList, nil
}

func (d *Derivative) computeDerivatives(timeSeries *TimeSeries) []float64 {
	zippedSeries := timeSeries.Zip()
	derivatives := make([]float64, 0, len(zippedSeries))

	for i, timestamp := range timeSeries.Timestamps {
		if i > 0 {
			preTimestamp := timeSeries.Timestamps[i-1]
			preValue := zippedSeries[preTimestamp]

			currentValue := zippedSeries[timestamp]
			delta := timestamp - preTimestamp

			derivative := 0.0
			if delta != 0 {
				derivative = (currentValue - preValue) / delta
			} else {
				derivative = currentValue - preValue
			}
			derivatives = append(derivatives, math.Abs(derivative))
		}
	}

	if len(derivatives) != 0 {
		derivatives = insertAt(derivatives, 0, derivatives[0])
	}
	return derivatives
}
