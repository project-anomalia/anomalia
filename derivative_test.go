package anomalia

import "testing"

func TestRunWithDerivative(t *testing.T) {
	timeSeries := &TimeSeries{
		Timestamps: []float64{1, 1, 3, 4, 5, 6, 7, 8, 9, 10},
		Values:     []float64{56, 59, 52, 49, 49, 1.5, 48, 50, 53, 44},
	}

	scoreList := NewDerivative().WithSmoothingFactor(0.3).Run(timeSeries)
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}

	if len(scoreList.Scores) != len(timeSeries.Values) {
		t.Fatalf("score list and time series dimensions do not match")
	}
}
