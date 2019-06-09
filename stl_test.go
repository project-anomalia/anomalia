package anomalia

import "testing"

func TestRunWithSTL(t *testing.T) {
	ts := NewTimeSeriesFromCSV("testdata/co2.csv")
	scoreList := NewSTL().WithWidth(35).WithPeriodicity(12).Run(ts)

	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}

	if len(scoreList.Scores) != ts.Size() {
		t.Fatalf("score list must have the same dimension as original time series")
	}
}
