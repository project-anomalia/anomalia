package anomalia

import (
	"math"
	"testing"
)

func TestRunSpearmanCorrelatorWhenTimeSeriesExactlyTheSame(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})

	coefficient := NewSpearmanCorrelator(timeSeriesA, timeSeriesB).Run()
	if coefficient != 1 {
		t.Fatalf("must return exactly 1")
	}
}

func TestRunSpearmanCorrelatorWhenTimeSeriesHaveNoLinearRelation(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{0, 3.2, 5.5, 7.1, 8.9, 9, 10.1, 10.5})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{-0.5, 1, 2.5, 4.1, 4.6, -1, 1, -1})

	coefficient := NewSpearmanCorrelator(timeSeriesA, timeSeriesB).Run()
	if math.Round(coefficient) != 0 {
		t.Fatalf("must return number close to 0")
	}
}

func TestRunSpearmanCorrelatorWhenTimeSeriesHaveDifferentSizes(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4}, []float64{0, 3.2, 5.5, 7.1, 8.9})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5}, []float64{-0.5, 1, 2.5, 4.1, 4.6, -1})

	// Assert panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("correlator did not panic")
		}
	}()

	NewSpearmanCorrelator(timeSeriesA, timeSeriesB).Run()
}

func TestRunSpearmanCorrelatorExample(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []float64{35, 23, 47, 17, 10, 43, 9, 6, 28})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []float64{30, 33, 45, 23, 8, 49, 12, 4, 31})

	coefficient := NewSpearmanCorrelator(timeSeriesA, timeSeriesB).Run()
	if coefficient != 0.9 {
		t.Fatalf("incorrect rank correlation coefficient")
	}
}
