package anomalia

import "testing"

func TestRunWithBitmap(t *testing.T) {
	timeSeries := &TimeSeries{
		Timestamps: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Values:     []float64{1, 5, 52, 49, 49, 1.5, 48, 50, 53, 44},
	}

	bitmap := &Bitmap{
		ChunkSize:        2,
		Precision:        4,
		LagWindowSize:    0,
		FutureWindowSize: 1,
	}

	scoreList := bitmap.Run(timeSeries)
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}

	if len(scoreList.Scores) != len(timeSeries.Timestamps) {
		t.Fatalf("both time series and score list dimensions do not match")
	}
}
