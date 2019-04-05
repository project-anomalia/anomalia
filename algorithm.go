package anomalia

// Algorithm is the base interface of all algorithms
type Algorithm interface {
	Run(*TimeSeries) *ScoreList
	computeScores(*TimeSeries) (*ScoreList, error)
}
