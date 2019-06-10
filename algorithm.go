package anomalia

// Algorithm is the base interface of all algorithms
type Algorithm interface {
	Run(*TimeSeries) *ScoreList
	computeScores(*TimeSeries) (*ScoreList, error)
}

// TimePeriod represents a time period marked by start and end timestamps.
type TimePeriod struct {
	Start float64
	End   float64
}
