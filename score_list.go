package anomalia

const noisePercentageThreshold = 0.001

// ScoreList holds timestamps and their scores
type ScoreList struct {
	Timestamps []float64
	Scores     []float64
}

// Denoise sets low(noisy) scores to 0.0
func (sl *ScoreList) Denoise() *ScoreList {
	threshold := noisePercentageThreshold * sl.Max()

	denoised := mapSlice(sl.Scores, func(score float64) float64 {
		if score < threshold {
			return 0.0
		}
		return score
	})
	return &ScoreList{sl.Timestamps, denoised}
}

// Max returns the maximum of the scores
func (sl *ScoreList) Max() float64 {
	_, max := minMax(sl.Scores)
	return max
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
