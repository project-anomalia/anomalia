package anomalia

const noisePercentageThreshold = 0.001

// ScoreList holds timestamps and their scores
type ScoreList struct {
	Timestamps []float64
	Scores     []float64
}

// Zip convert the score list to map (map[Timestamp]Score)
func (sl *ScoreList) Zip() map[float64]float64 {
	m := make(map[float64]float64)
	sorted := sortedCopy(sl.Timestamps)

	for idx, timestamp := range sorted {
		m[timestamp] = sl.Scores[idx]
	}
	return m
}

// Denoise sets low(noisy) scores to 0.0
func (sl *ScoreList) Denoise() *ScoreList {
	_, max := minMax(sl.Scores)
	threshold := noisePercentageThreshold * max

	denoised := mapSlice(sl.Scores, func(score float64) float64 {
		if score < threshold {
			return 0.0
		} else {
			return score
		}
	})
	return &ScoreList{sl.Timestamps, denoised}
}
