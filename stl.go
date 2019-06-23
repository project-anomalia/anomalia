package anomalia

import "github.com/project-anomalia/stl"

type STLMethod int32

const (
	// Additive method suggests that the components are added together (linear model).
	Additive STLMethod = iota

	// Multiplicative method suggests that the components are multiplied together (non-linear model).
	Multiplicative
)

// STL holds Seasonal-Trend With Loess algorithm configuration.
//
// The STL algorithm decomposes a time series into seasonal, trend and remainder components.
// The paper describing this algorithm can found here: https://search.proquest.com/openview/cc5001e8a0978a6c029ae9a41af00f21
type STL struct {
	periodicity         int
	width               int
	robustIterations    stl.Opt
	iterations          stl.Opt
	seasonalConfig      *stl.Config
	trendConfig         *stl.Config
	lowPassFilterConfig *stl.Config
	method              stl.ModelType
}

// NewSTL returns an instance of the STL struct.
func NewSTL() *STL {
	return &STL{
		robustIterations: stl.WithRobustIter(0),
		iterations:       stl.WithIter(2),
		method:           stl.Additive(),
	}
}

func (s *STL) Periodicity(p int) *STL {
	s.periodicity = p
	return s
}

func (s *STL) Width(w int) *STL {
	s.width = w
	return s
}

func (s *STL) RobustIterations(n int) *STL {
	s.robustIterations = stl.WithRobustIter(n)
	return s
}

func (s *STL) Iterations(n int) *STL {
	s.iterations = stl.WithIter(n)
	return s
}

func (s *STL) SeasonalConfig(config *stl.Config) *STL {
	s.seasonalConfig = config
	return s
}

func (s *STL) TrendConfig(config *stl.Config) *STL {
	s.trendConfig = config
	return s
}

func (s *STL) LowPassFilterConfig(config *stl.Config) *STL {
	s.lowPassFilterConfig = config
	return s
}

func (s *STL) MethodType(method STLMethod) *STL {
	switch method {
	case Additive:
		s.method = stl.Additive()
	case Multiplicative:
		s.method = stl.Multiplicative()
	default:
		panic("invalid STL method type")
	}
	return s
}

// Run runs the STL algorithm over the time series.
func (s *STL) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := s.computeScores(timeSeries)
	return scoreList
}

func (s *STL) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	options := []stl.Opt{s.iterations, s.robustIterations}

	if s.seasonalConfig != nil {
		options = append(options, stl.WithSeasonalConfig(*s.seasonalConfig))
	}

	if s.trendConfig != nil {
		options = append(options, stl.WithTrendConfig(*s.trendConfig))
	}

	if s.lowPassFilterConfig != nil {
		options = append(options, stl.WithLowpassConfig(*s.lowPassFilterConfig))
	}

	result := stl.Decompose(timeSeries.Values, s.periodicity, s.width, s.method, options...)
	if result.Err != nil {
		return nil, result.Err
	}
	return &ScoreList{timeSeries.Timestamps, result.Resid}, nil
}
