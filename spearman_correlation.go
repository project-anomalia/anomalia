package anomalia

import "sort"

// SpearmanCorrelation holds the Spearman Correlation algorithm configuration.
// It is the non-parametric version of the Pearson correlation and it should be used
// when the time series distribution is unknown or not normally distributed.
//
// Spearman’s correlator returns a value from -1 to 1, where:
//	- +1  = a perfect positive correlation between ranks
//	- -1  = a perfect negative correlation between ranks
//	- 0   = no correlation between ranks.
type SpearmanCorrelation struct {
	current, target *TimeSeries
}

type rank struct{ x, y, xRank, yRank float64 }

// NewSpearmanCorrelation returns an instance of the spearman correlation struct.
func NewSpearmanCorrelation(current, target *TimeSeries) *SpearmanCorrelation {
	return &SpearmanCorrelation{current, target}
}

// Run runs the spearman correlation on the current and target time series.
func (sc *SpearmanCorrelation) Run() float64 {
	sc.sanityCheck()

	// Build up the ranks slice
	ranks := make([]rank, sc.current.Size())
	for index, currentValue := range sc.current.Values {
		ranks[index] = rank{x: currentValue, y: sc.target.Values[index]}
	}

	// Sort the ranks by x
	sort.Slice(ranks, func(i, j int) bool { return ranks[i].x < ranks[j].x })

	// Rank the current series
	for pos := 0; pos < len(ranks); pos++ {
		ranks[pos].xRank = float64(pos) + 1
		duplicateValues := []int{pos}
		for nested, p := range ranks {
			if ranks[pos].x == p.x {
				if pos != nested {
					duplicateValues = append(duplicateValues, nested)
				}
			}
		}

		sum := 0
		for _, val := range duplicateValues {
			sum += val
		}

		avg := float64(sum+len(duplicateValues)) / float64(len(duplicateValues))
		ranks[pos].xRank = avg
		for index := 1; index < len(duplicateValues); index++ {
			ranks[duplicateValues[index]].xRank = avg
		}
		pos += len(duplicateValues) - 1
	}

	// Sort the ranks by y
	sort.Slice(ranks, func(i int, j int) bool { return ranks[i].y < ranks[j].y })

	// Rank the target series
	for pos := 0; pos < len(ranks); pos++ {
		ranks[pos].yRank = float64(pos) + 1
		duplicateValues := []int{pos}
		for nested, p := range ranks {
			if ranks[pos].y == p.y {
				if pos != nested {
					duplicateValues = append(duplicateValues, nested)
				}
			}
		}

		sum := 0
		for _, val := range duplicateValues {
			sum += val
		}

		avg := float64(sum+len(duplicateValues)) / float64(len(duplicateValues))
		ranks[pos].yRank = avg
		for index := 1; index < len(duplicateValues); index++ {
			ranks[duplicateValues[index]].yRank = avg
		}
		pos += len(duplicateValues) - 1
	}

	// Adapt both current and target series
	for index, rank := range ranks {
		sc.current.Values[index] = rank.xRank
		sc.target.Values[index] = rank.yRank
	}

	return NewPearsonCorrelation(sc.current, sc.target).Run()
}

func (sc *SpearmanCorrelation) sanityCheck() {
	if sc.current.Size() < 3 || sc.current.Size() != sc.target.Size() {
		panic("current and/or target series have an invalid dimension")
	}
}
