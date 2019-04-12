package anomalia

import "testing"

func TestRunWithWeightedSum(t *testing.T) {
	weightedSum := &WeightedSum{
		EmaWeight:                0.65,
		EmaSignificant:           0.94,
		ExponentialMovingAverage: &ExponentialMovingAverage{2, 0.2},
		Derivative:               &Derivative{0.2},
	}

	timeSeries := &TimeSeries{
		Timestamps: []float64{1, 1, 3, 4, 5, 6, 7, 8, 9, 10},
		Values:     []float64{56, 59, 52, 49, 49, 1.5, 48, 50, 53, 44},
	}

	scoreList := weightedSum.Run(timeSeries)
	if len(scoreList.Timestamps) != len(timeSeries.Timestamps) {
		t.Fatalf("score list and time series dimensions do not match")
	}
}
