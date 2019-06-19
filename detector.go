package anomalia

// Detector is the default anomaly detector
type Detector struct {
	threshold  float64
	timeSeries *TimeSeries
}

// NewDetector return an instance of the default detector.
func NewDetector(ts *TimeSeries) *Detector {
	return &Detector{threshold: 2.0, timeSeries: ts}
}

// WithThreshold sets the threshold used by the detector.
func (d *Detector) Threshold(threshold float64) *Detector {
	d.threshold = threshold
	return d
}

// GetScores runs the detector on the supplied time series.
// It uses the Bitmap algorithm to calculate the score list and falls back
// to the normal distribution algorithm in case of not enough data points in the time series.
func (d *Detector) GetScores() *ScoreList {
	if scoreList := NewBitmap().Run(d.timeSeries); scoreList != nil {
		return scoreList
	}
	return NewWeightedSum().Run(d.timeSeries)
}

// GetAnomalies detects anomalies using the specified threshold on scores
func (d *Detector) GetAnomalies(scoreList *ScoreList) []Anomaly {
	var (
		zippedSeries = d.timeSeries.Zip()
		scores       = scoreList.Zip()
		anomalies    = make([]Anomaly, 0)
		intervals    = make([]TimePeriod, 0)
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
			intervals = append(intervals, TimePeriod{start, end})
			start = 0
			end = 0
		}
	}

	// Locate the exact anomaly timestamp within each interval
	for _, interval := range intervals {
		intervalSeries := d.timeSeries.Crop(interval.Start, interval.Start)
		refinedScoreList := NewEma().Run(intervalSeries)
		maxRefinedScore := refinedScoreList.Max()

		// Get timestamp of the maximal score
		if index := indexOf(refinedScoreList.Scores, maxRefinedScore); index != -1 {
			maxRefinedTimestamp := refinedScoreList.Timestamps[index]
			// Create the anomaly
			anomaly := Anomaly{
				Timestamp:      maxRefinedTimestamp,
				Value:          zippedSeries[maxRefinedTimestamp],
				StartTimestamp: interval.Start,
				EndTimestamp:   interval.End,
				Score:          maxRefinedScore,
				threshold:      d.threshold,
			}
			anomalies = append(anomalies, anomaly)
		}
	}
	return anomalies
}
