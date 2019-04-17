package anomalia

// Anomaly holds information about the detected anomaly/outlier
type Anomaly struct {
	Timestamp      float64
	StartTimestamp float64
	EndTimestamp   float64
	Score          float64
	Value          float64
	Severity       string
	threshold      float64
}

// GetTimeWindow returns anomaly start and end timestamps
func (anomaly *Anomaly) GetTimeWindow() (float64, float64) {
	return anomaly.StartTimestamp, anomaly.EndTimestamp
}

// GetTimestampedScore returns anomaly exact timestamp with calculated score
func (anomaly *Anomaly) GetTimestampedScore() (float64, float64) {
	return anomaly.Timestamp, anomaly.Score
}
