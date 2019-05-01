package anomalia

import "testing"

func TestNewCorrelator(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{1, 2, 3, 4}, []float64{0.5, 1, 0.2, 5})
	timeSeriesB := NewTimeSeries([]float64{4, 6, 7, 8}, []float64{1.5, 3.2, 1.2, 4})
	correlator := NewCorrelator(timeSeriesA, timeSeriesB).
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
	result1 := NewCorrelator(timeSeriesA, timeSeriesB).Run()
	result2 := NewCorrelator(timeSeriesA, timeSeriesC).Run()

	if result1.Coefficient != result2.Coefficient {
		t.Fatalf("correlation coefficient did not match")
	}

	if result1.Shift != result2.Shift {
		t.Fatalf("correlation shift did not match")
	}
}
