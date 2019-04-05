package anomalia

import (
	"testing"
)

func TestRunWithNormalDistribution(t *testing.T) {
	timeSeries := &TimeSeries{
		Timestamps: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Values:     []float64{56, 59, 52, 49, 49, 1.5, 48, 50, 53, 44},
	}

	scoreList := (&NormalDistribution{0.0025}).Run(timeSeries)
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}

	anomalies := filter(scoreList.Scores, func(val float64) bool { return val != 0.0 })
	if len(anomalies) != 1 {
		t.Fatalf("only 1 anomaly in the data set")
	}
}
