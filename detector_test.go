package anomalia

import "testing"

func TestRunDefaultDetectorOnSmallDataset(t *testing.T) {
	timeSeries := generateFakeTimeSeries(100)
	scoreList := NewDetector(timeSeries).GetScores()
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}
}

func TestRunDefaultDetectorOnLargeDataset(t *testing.T) {
	timeSeries := generateFakeTimeSeries(3000)
	scoreList := NewDetector(timeSeries).Threshold(4.5).GetScores()
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}
}

func TestGetAnomaliesUsingDefaultDetector(t *testing.T) {
	timeSeries := generateFakeTimeSeries(2000)
	detector := NewDetector(timeSeries).Threshold(3.0)

	scores := detector.GetScores()
	anomalies := detector.GetAnomalies(scores)
	if len(anomalies) != 1 {
		t.Fatalf("should be a least one anomaly")
	}
}

func TestGetAnomaliesInTestData(t *testing.T) {
	ts := NewTimeSeriesFromCSV("testdata/airline-passengers.csv")
	detector := NewDetector(ts).Threshold(1.3)

	scoreList := NewSTL().WithWidth(15).WithPeriodicity(12).WithMethod(Multiplicative).Run(ts).Denoise()
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}

	anomalies := detector.GetAnomalies(scoreList)
	if len(anomalies) != 2 {
		t.Fatalf("there are exactly 2 anomalies")
	}
}
