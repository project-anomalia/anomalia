package anomalia

// ScoreList holds timestamps and their scores
type ScoreList struct {
	Timestamps []float64
	Scores     []float64
}

// Algorithm is the base interface of all algorithms
type Algorithm interface {
	Run(*TimeSeries) *ScoreList
	computeScores(*TimeSeries) (*ScoreList, error)
}
