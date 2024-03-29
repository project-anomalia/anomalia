package anomalia

import "testing"

func TestDenoiseScoreList(t *testing.T) {
	denoised := fakeScoreList().Denoise()
	if denoised == nil {
		t.Fatalf("score list cannot be nil")
	}

	noisyScore := denoised.Scores[0]
	if noisyScore != 0.0 {
		t.Fatalf("noisy scores should be zeroed")
	}
}

func TestMaxOfScoreList(t *testing.T) {
	maxScore := fakeScoreList().Max()
	if maxScore != 4.6 {
		t.Fatalf("max score is incorrect")
	}
}

func fakeScoreList() *ScoreList {
	return &ScoreList{
		Timestamps: []float64{1, 2, 3, 4, 5, 6},
		Scores:     []float64{0.0010, 4.6, 4.6, 4.6, 1.0, 1.0},
	}
}
