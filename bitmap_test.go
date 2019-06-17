package anomalia

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	mu        sync.Mutex
	generator = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func TestRunWithBitmap(t *testing.T) {
	timeSeries := generateFakeTimeSeries(2000)
	scoreList := NewBitmap().ChunkSize(3).Precision(5).Run(timeSeries)
	if scoreList == nil {
		t.Fatalf("score list cannot be nil")
	}

	if len(scoreList.Scores) != len(timeSeries.Timestamps) {
		t.Fatalf("both time series and score list dimensions do not match")
	}
}

func TestRunBitmapWhenNotEnoughDataPoints(t *testing.T) {
	timeSeries := &TimeSeries{
		Timestamps: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Values:     []float64{1, 5, 52, 49, 49, 1.5, 48, 50, 53, 44},
	}
	scoreList := NewBitmap().Run(timeSeries)
	if scoreList != nil {
		t.Fatalf("score list must be nil (not enough data points)")
	}
}

func generateFakeTimeSeries(datasetSize int) *TimeSeries {
	timestamps := make([]float64, datasetSize)
	for i := 0; i < datasetSize; i++ {
		timestamps[i] = float64(i) + 1
	}
	mu.Lock()
	values := make([]float64, datasetSize)
	for i := 0; i < datasetSize; i++ {
		values[i] = RandomSineValue(generator, datasetSize)
	}
	defer mu.Unlock()
	return &TimeSeries{timestamps, values}
}
