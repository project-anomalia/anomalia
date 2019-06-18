package anomalia

import "testing"

func TestRunCorrelatorWithXCorr(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})

	coefficient := NewCorrelator(timeSeriesA, timeSeriesB).
		CorrelationMethod(XCorr, []float64{30, 0.01}).
		UseAnomalyScore(true).
		Run()
	if coefficient != 1.0 {
		t.Fatalf("incorrect coefficient: time series are exactly the same")
	}
}

func TestRunCorrelatorWithSpearmanRank(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})

	coefficient := NewCorrelator(timeSeriesA, timeSeriesB).
		CorrelationMethod(SpearmanRank, nil).
		TimePeriod(0, 2).
		Run()
	if coefficient != 1.0 {
		t.Fatalf("incorrect coefficient: time series are exactly the same")
	}
}

func TestRunCorrelatorWithPearson(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})

	coefficient := NewCorrelator(timeSeriesA, timeSeriesB).CorrelationMethod(Pearson, nil).Run()
	if coefficient != 1.0 {
		t.Fatalf("incorrect coefficient: time series are exactly the same")
	}
}

func TestRunPearsonCorrelationWhenTimeSeriesHaveDifferentSizes(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4}, []float64{0, 3.2, 5.5, 7.1, 8.9})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5}, []float64{-0.5, 1, 2.5, 4.1, 4.6, -1})

	// Assert panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("correlator did not panic")
		}
	}()

	NewCorrelator(timeSeriesA, timeSeriesB).CorrelationMethod(Pearson, nil).Run()
}

func TestXCorrelationWhenNotEnoughDataPoints(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1}, []float64{0.5, 0})
	timeSeriesB := NewTimeSeries([]float64{0}, []float64{0.5})

	// Assert panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("correlator did not panic")
		}
	}()

	NewCorrelator(timeSeriesA, timeSeriesB).UseAnomalyScore(true).Run()
}

func TestRunSpearmanCorrelationWhenTimeSeriesHaveDifferentSizes(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4}, []float64{0, 3.2, 5.5, 7.1, 8.9})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5}, []float64{-0.5, 1, 2.5, 4.1, 4.6, -1})

	// Assert panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("correlator did not panic")
		}
	}()

	NewCorrelator(timeSeriesA, timeSeriesB).CorrelationMethod(SpearmanRank, nil).Run()
}
