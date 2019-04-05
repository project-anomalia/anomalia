package anomalia

import (
	"testing"
)

func TestDenoiseScoreList(t *testing.T) {
	scoreList := &ScoreList{
		Timestamps: []float64{1, 2, 3, 4, 5, 6},
		Scores:     []float64{0.0010, 4.6, 4.6, 4.6, 1.0, 1.0},
	}
	denoised := scoreList.Denoise()
	if denoised == nil {
		t.Fatalf("score list cannot be nil")
	}

	noisyScore := denoised.Scores[0]
	if noisyScore != 0.0 {
		t.Fatalf("noisy scores should be zeroed")
	}
}
