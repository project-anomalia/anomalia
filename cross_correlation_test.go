package anomalia

import "testing"

func TestNewCrossCorrelation(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{2, 3, -2, 3, 2, 4, 1, -1})
	correlator := NewCrossCorrelation(timeSeriesA, timeSeriesB).
		WithMaxShift(30).
		WithImpact(0.01)
	if correlator == nil {
		t.Fatalf("failed to initialize correlator")
	}
}

func TestRunCrossCorrelation(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []float64{0, 0, 0, 0, 0.5, 1, 1, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []float64{0, 0.5, 1, 1, 1, 0, 0, 0, 0})
	timeSeriesC := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5}, []float64{0, 0.5, 1, 1, 1, 0})
	result1 := NewCrossCorrelation(timeSeriesA, timeSeriesB).GetCorrelationResult()
	result2 := NewCrossCorrelation(timeSeriesA, timeSeriesC).GetCorrelationResult()

	if result1.Coefficient != result2.Coefficient {
		t.Fatalf("correlation coefficient did not match")
	}

	if result1.Shift != result2.Shift {
		t.Fatalf("correlation shift did not match")
	}
}

func TestCorrelationWhenNotEnoughDataPoints(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1}, []float64{0.5, 0})
	timeSeriesB := NewTimeSeries([]float64{0}, []float64{0.5})

	// Assert panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("correlator did not panic")
		}
	}()
	NewCrossCorrelation(timeSeriesA, timeSeriesB).Run()
}

func TestCorrelationWhenTimeSeriesExactlyTheSame(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	result := NewCrossCorrelation(timeSeriesA, timeSeriesB).Run()
	if result != 1 {
		t.Fatalf("incorrect coefficient: time series are exactly the same")
	}
}
