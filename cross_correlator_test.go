package anomalia

import "testing"

func TestNewCorrelator(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{2, 3, -2, 3, 2, 4, 1, -1})
	correlator := NewCrossCorrelator(timeSeriesA, timeSeriesB).
		WithMaxShift(30).
		WithImpact(0.01).
		UseAnomalyScore(true)
	if correlator == nil {
		t.Fatalf("failed to initialize correlator")
	}
}

func TestRunCrossCorrelator(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []float64{0, 0, 0, 0, 0.5, 1, 1, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []float64{0, 0.5, 1, 1, 1, 0, 0, 0, 0})
	timeSeriesC := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5}, []float64{0, 0.5, 1, 1, 1, 0})
	result1 := NewCrossCorrelator(timeSeriesA, timeSeriesB).Run()
	result2 := NewCrossCorrelator(timeSeriesA, timeSeriesC).Run()

	if result1.Coefficient != result2.Coefficient {
		t.Fatalf("correlation coefficient did not match")
	}

	if result1.Shift != result2.Shift {
		t.Fatalf("correlation shift did not match")
	}
}

func TestCorrelatorWhenNotEnoughDataPoints(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1}, []float64{0.5, 0})
	timeSeriesB := NewTimeSeries([]float64{0}, []float64{0.5})

	// Assert panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("correlator did not panic")
		}
	}()
	NewCrossCorrelator(timeSeriesA, timeSeriesB).Run()
}

func TestCorrelatorWhenTimeSeriesExactlyTheSame(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	result := NewCrossCorrelator(timeSeriesA, timeSeriesB).Run()
	if result.Coefficient != 1 {
		t.Fatalf("incorrect coefficient: time series are exactly the same")
	}
}
