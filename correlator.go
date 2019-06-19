package anomalia

// CorrelationAlgorithm base interface for correlation algorithms.
type CorrelationAlgorithm interface {
	Run() float64
	sanityCheck() error
}

// CorrelationMethod type checker for correlation method
type CorrelationMethod int32

const (
	// XCorr represents the Cross Correlation algorithm.
	XCorr CorrelationMethod = iota
	// SpearmanRank represents the Spearman Rank Correlation algorithm.
	SpearmanRank
	// Pearson represents the Pearson Correlation algorithm.
	Pearson
)

// Correlator holds the correlator configuration.
type Correlator struct {
	current, target *TimeSeries
	algorithm       CorrelationAlgorithm
	useAnomalyScore bool
}

// NewCorrelator returns an instance of the correlation algorithm.
func NewCorrelator(current, target *TimeSeries) *Correlator {
	if current == nil || target == nil {
		panic("either current or target time series cannot be nil")
	}
	return &Correlator{
		current: current,
		target:  target,
	}
}

// CorrelationMethod specifies which correlation method to use (XCross or SpearmanRank).
func (c *Correlator) CorrelationMethod(method CorrelationMethod, options []float64) *Correlator {
	c.algorithm = c.getCorrelationAlgorithmByMethod(method, options)
	return c
}

// WithTimePeriod crops the current and target time series to specified range.
func (c *Correlator) TimePeriod(start, end float64) *Correlator {
	c.current = c.current.Crop(start, end)
	c.target = c.target.Crop(start, end)
	return c
}

// UseAnomalyScore tells the correlator to calculate anomaly scores from both time series.
func (c *Correlator) UseAnomalyScore(use bool) *Correlator {
	c.useAnomalyScore = use
	return c
}

// Run runs the correlator.
func (c *Correlator) Run() float64 {
	if err := c.algorithm.sanityCheck(); err != nil {
		panic(err)
	}

	if c.useAnomalyScore {
		c.current = getAnomalyScores(NewDetector(c.current))
		c.target = getAnomalyScores(NewDetector(c.target))
	}

	return c.algorithm.Run()
}

func (c *Correlator) getCorrelationAlgorithmByMethod(method CorrelationMethod, options []float64) CorrelationAlgorithm {
	var algorithm CorrelationAlgorithm
	switch method {
	case XCorr:
		algorithm = NewCrossCorrelation(c.current, c.target)
		if options != nil && len(options) > 0 {
			algorithm = algorithm.(*CrossCorrelation).
				WithMaxShift(options[0]).
				WithImpact(options[1])
		}
	case SpearmanRank:
		algorithm = NewSpearmanCorrelation(c.current, c.target)
	case Pearson:
		algorithm = NewPearsonCorrelation(c.current, c.target)
	default:
		panic("unsupported correlation method/algorithm")
	}
	return algorithm
}

func getAnomalyScores(detector *Detector) *TimeSeries {
	scoreList := detector.GetScores()
	if scoreList == nil {
		panic("failed to calculate anomaly scores")
	}
	return &TimeSeries{scoreList.Timestamps, scoreList.Scores}
}
