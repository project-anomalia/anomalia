package anomalia

import "testing"

func TestRunWithAbsoluteThreshold(t *testing.T) {
	timeSeries := &TimeSeries{
		Timestamps: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Values:     []float64{0.4, 5.0, 5.0, 9.0, 1.0, 4.5, 3.0, 6.0, 1.4, 5.3},
	}
	scoreList := NewAbsoluteThreshold().WithBounds(0.5, 2.5).Run(timeSeries)
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}

	if len(scoreList.Scores) != len(timeSeries.Values) {
		t.Fatalf("score list and time series dimensions do not match")
	}
}
