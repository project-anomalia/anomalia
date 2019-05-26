package anomalia

import "testing"

func TestRunDefaultDetectorOnSmallDataset(t *testing.T) {
	timeSeries := generateFakeTimeSeries(100)
	scoreList := NewDetector().GetScores(timeSeries)
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}
}

func TestRunDefaultDetectorOnLargeDataset(t *testing.T) {
	timeSeries := generateFakeTimeSeries(3000)
	scoreList := NewDetector().WithThreshold(4.5).GetScores(timeSeries)
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}
}

func TestGetAnomaliesUsingDefaultDetector(t *testing.T) {
	timeSeries := generateFakeTimeSeries(2000)
	anomalies := NewDetector().WithThreshold(3.0).GetAnomalies(timeSeries)
	if len(anomalies) != 1 {
		t.Fatalf("should be a least one anomaly")
	}
}
