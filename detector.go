package anomalia

// Detector is the default anomaly detector
type Detector struct {
	threshold float64
}

type tuple struct{ start, end float64 }

// NewDetector return an instance of the default detector.
func NewDetector() *Detector {
	return &Detector{2.0}
}

// WithThreshold sets the threshold used by the detector.
func (d *Detector) WithThreshold(threshold float64) *Detector {
	d.threshold = threshold
	return d
}

// GetScores runs the detector on the supplied time series
// It uses the bitmap alogrithm to calculate the score list and falls back
// to the normal distribution algorithm in case of not enough data points in the time series.
func (d *Detector) GetScores(timeSeries *TimeSeries) *ScoreList {
	if scoreList := NewBitmap().Run(timeSeries); scoreList != nil {
		return scoreList
	}
	return NewWeightedSum().Run(timeSeries)
}

// GetAnomalies detects anomalies using the specified threshold on scores
func (d *Detector) GetAnomalies(timeSeries *TimeSeries) []Anomaly {
	var (
		zippedSeries = timeSeries.Zip()
		scoreList    = d.GetScores(timeSeries)
		scores       = scoreList.Zip()
		anomalies    = make([]Anomaly, 0)
		intervals    = make([]tuple, 0)
	)

	// Find all anomalies intervals
	var start, end float64
	for _, timestamp := range scoreList.Timestamps {
		if scores[timestamp] > d.threshold {
			end = timestamp
			if start == 0 {
				start = timestamp
			}
		} else if (start != 0) && (end != 0) {
			intervals = append(intervals, tuple{start, end})
			start = 0
			end = 0
		}
	}

	// Locate the exact anomaly timestamp within each interval
	for _, interval := range intervals {
		intervalSeries := timeSeries.Crop(interval.start, interval.end)
		refinedScoreList := NewEma().Run(intervalSeries)
		maxRefinedScore := refinedScoreList.Max()

		// Get timestamp of the maximal score
		if index := indexOf(refinedScoreList.Scores, maxRefinedScore); index != -1 {
			maxRefinedTimestamp := refinedScoreList.Timestamps[index]
			// Create the anomaly
			anomaly := Anomaly{
				Timestamp:      maxRefinedTimestamp,
				Value:          zippedSeries[maxRefinedTimestamp],
				StartTimestamp: interval.start,
				EndTimestamp:   interval.end,
				Score:          maxRefinedScore,
				threshold:      d.threshold,
			}
			anomalies = append(anomalies, anomaly)
		}
	}
	return anomalies
}
