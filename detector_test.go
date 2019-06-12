package anomalia

import "testing"

func TestRunDefaultDetectorOnSmallDataset(t *testing.T) {
	timeSeries := generateFakeTimeSeries(100)
	scoreList := NewDetector().WithTimeSeries(timeSeries).GetScores()
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}
}

func TestRunDefaultDetectorOnLargeDataset(t *testing.T) {
	timeSeries := generateFakeTimeSeries(3000)
	scoreList := NewDetector().
		WithTimeSeries(timeSeries).
		WithThreshold(4.5).
		GetScores()
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}
}

func TestGetAnomaliesUsingDefaultDetector(t *testing.T) {
	timeSeries := generateFakeTimeSeries(2000)
	detector := NewDetector().WithTimeSeries(timeSeries).WithThreshold(3.0)

	scores := detector.GetScores()
	anomalies := detector.GetAnomalies(scores)
	if len(anomalies) != 1 {
		t.Fatalf("should be a least one anomaly")
	}
}
