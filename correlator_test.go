package anomalia

import "testing"

func TestRunCorrelatorWithXCorr(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})

	coefficient := NewCorrelator().WithTimeSeries(timeSeriesA, timeSeriesB).
		WithMethod(XCorr, []float64{30, 0.01}).
		UseAnomalyScore(true).
		Run()
	if coefficient != 1.0 {
		t.Fatalf("incorrect coefficient: time series are exactly the same")
	}
}

func TestRunCorrelatorWithSpearmanRank(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})

	coefficient := NewCorrelator().WithTimeSeries(timeSeriesA, timeSeriesB).
		WithMethod(SpearmanRank, nil).
		WithTimePeriod(0, 2).
		Run()
	if coefficient != 1.0 {
		t.Fatalf("incorrect coefficient: time series are exactly the same")
	}
}
