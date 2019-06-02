package anomalia

import (
	"math"
	"testing"
)

func TestRunPearsonCorrelationWhenTimeSeriesExactlyTheSame(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{1, 2, -2, 4, 2, 3, 1, 0})

	coefficient := NewPearsonCorrelation(timeSeriesA, timeSeriesB).Run()
	if coefficient != 1 {
		t.Fatalf("must return exactly 1")
	}
}

func TestRunPearsonCorrelationWhenTimeSeriesHaveNoLinearRelation(t *testing.T) {
	timeSeriesA := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{0, 3.2, 5.5, 7.1, 8.9, 9, 10.1, 10.5})
	timeSeriesB := NewTimeSeries([]float64{0, 1, 2, 3, 4, 5, 6, 7}, []float64{-0.5, 1, 2.5, 4.1, 4.6, -1, 1, -1})

	coefficient := NewPearsonCorrelation(timeSeriesA, timeSeriesB).Run()
	if math.Round(coefficient) != 0 {
		t.Fatalf("must return number close to 0")
	}
}
