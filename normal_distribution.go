package anomalia

const defaultEpsilonThreshold = 0.0025

// NormalDistribution holds the normal distribution algorithm configuration.
type NormalDistribution struct {
	epsilonThreshold float64
}

// NewNormalDistribution returns normal distribution instance.
func NewNormalDistribution() *NormalDistribution {
	return &NormalDistribution{defaultEpsilonThreshold}
}

// EpsilonThreshold sets the Gaussian epsilon threshold.
func (nd *NormalDistribution) EpsilonThreshold(threshold float64) *NormalDistribution {
	nd.epsilonThreshold = threshold
	return nd
}

// Run runs the normal distribution algorithm over the time series.
func (nd *NormalDistribution) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := nd.computeScores(timeSeries)
	return scoreList
}

func (nd *NormalDistribution) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	mean := timeSeries.Average()
	std := timeSeries.Stdev()

	scores := mapSlice(timeSeries.Values, func(value float64) float64 {
		score := Pdf(mean, std)(value)
		if score < nd.epsilonThreshold {
			return score
		}
		return 0.0
	})

	scoreList := &ScoreList{timeSeries.Timestamps, scores}
	return scoreList, nil
}
