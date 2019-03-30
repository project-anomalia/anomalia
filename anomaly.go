package anomalia

// Anomaly holds information about the detected anomaly/outlier
type Anomaly struct {
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
