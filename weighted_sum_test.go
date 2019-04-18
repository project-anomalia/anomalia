package anomalia

import "testing"

func TestRunWithWeightedSum(t *testing.T) {
	timeSeries := &TimeSeries{
		Timestamps: []float64{1, 1, 3, 4, 5, 6, 7, 8, 9, 10},
		Values:     []float64{56, 59, 52, 49, 49, 1.5, 48, 50, 53, 44},
	}

	scoreList := NewWeightedSum().Run(timeSeries)
	if len(scoreList.Timestamps) != len(timeSeries.Timestamps) {
		t.Fatalf("score list and time series dimensions do not match")
	}
}
