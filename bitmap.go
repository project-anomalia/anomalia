package anomalia

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const minimalPointsInWindows = 50

// Bitmap holds bitmap algorithm configuration.
//
// The Bitmap algorithm breaks the time series into chunks and uses
// the frequency of similar chunks to determine anomalies scores.
// The scoring happens by sliding both lagging and future windows.
type Bitmap struct {
	chunkSize        int
	precision        int
	lagWindowSize    int
	futureWindowSize int
}

// NewBitmap returns Bitmap instance.
func NewBitmap() *Bitmap {
	return &Bitmap{
		chunkSize:        2,
		precision:        4,
		lagWindowSize:    0,
		futureWindowSize: 0,
	}
}

// ChunkSize sets the chunk size to use (defaults to 2).
func (b *Bitmap) ChunkSize(size int) *Bitmap {
	b.chunkSize = size
	return b
}

// Precision sets the precision.
func (b *Bitmap) Precision(p int) *Bitmap {
	b.precision = p
	return b
}

// LagWindowSize sets the lag window size (defaults to 0).
func (b *Bitmap) LagWindowSize(size int) *Bitmap {
	b.lagWindowSize = size
	return b
}

// FutureWindowSize sets the future window size (default to 0).
func (b *Bitmap) FutureWindowSize(size int) *Bitmap {
	b.futureWindowSize = size
	return b
}

// Run runs the bitmap algorithm over the time series
func (b *Bitmap) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := b.computeScores(timeSeries)
	return scoreList
}

func (b *Bitmap) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	// Update both lagging and future windows size
	b.lagWindowSize = int(0.0125 * float64(len(timeSeries.Timestamps)))
	b.futureWindowSize = int(0.0125 * float64(len(timeSeries.Timestamps)))

	// Perform sanity check
	if _, err := b.sanityCheck(timeSeries); err != nil {
		return nil, err
	}

	sax := b.generateSAX(timeSeries)
	laggingsMaps, futureMaps := b.constructAllSAXChunks(timeSeries, sax)
	dimension := timeSeries.Size()

	computeScoreBetweenTwoWindows := func(idx int) float64 {
		lagWindowChunk := laggingsMaps[idx]
		futureWindowChunk := futureMaps[idx]
		score := 0.0

		// Iterate over lagging window chunks
		for chunk := range lagWindowChunk {
			if _, ok := futureWindowChunk[chunk]; ok {
				score += math.Pow(float64(futureWindowChunk[chunk]-lagWindowChunk[chunk]), 2.0)
			} else {
				score += math.Pow(float64(lagWindowChunk[chunk]), 2.0)
			}
		}

		// Iterate over future window chunks
		for chunk := range futureWindowChunk {
			if _, ok := lagWindowChunk[chunk]; !ok {
				score += math.Pow(float64(futureWindowChunk[chunk]), 2)
			}
		}
		return score
	}

	scores := mapSliceWithIndex(timeSeries.Timestamps, func(idx int, timestamp float64) float64 {
		if (idx < b.lagWindowSize) || (idx > (dimension - b.futureWindowSize)) {
			return 0.0
		}
		return computeScoreBetweenTwoWindows(idx)
	})

	scoreList := &ScoreList{timeSeries.Timestamps, scores}
	return scoreList, nil
}

// generateSAX generates the SAX representation of the time series values
func (b *Bitmap) generateSAX(timeSeries *TimeSeries) BitmapBinary {
	sections := make(map[int]float64)
	min, max := minMax(timeSeries.Values)

	// Break the whole value range into different sections
	sectionHeight := (max - min) / float64(b.precision)
	for i := 0; i < b.precision; i++ {
		sections[i] = min + float64(i)*sectionHeight
	}

	// Generate SAX representation
	sectionsNumbers := mapIntKeys(sections)
	generateSingleSAX := func(value float64) string {
		sax := 0
		for _, sectionNumber := range sectionsNumbers {
			if value >= sections[sectionNumber] {
				sax = sectionNumber
			} else {
				break
			}
		}
		return strconv.Itoa(sax)
	}

	var saxBuilder strings.Builder
	for _, value := range timeSeries.Values {
		singleSAX := generateSingleSAX(value)
		saxBuilder.WriteString(singleSAX)
	}
	return BitmapBinary(saxBuilder.String())
}

func (b *Bitmap) constructChunkFrequencyMap(sax BitmapBinary) map[BitmapBinary]int {
	frequencyMap := make(map[BitmapBinary]int)
	saxLength := sax.Len()
	for i := 0; i < saxLength; i++ {
		if i+b.chunkSize <= saxLength {
			chunk := sax.Slice(i, i+b.chunkSize)
			frequencyMap[chunk]++
		}
	}
	return frequencyMap
}

func (b *Bitmap) constructAllSAXChunks(timeSeries *TimeSeries, sax BitmapBinary) (map[int]map[BitmapBinary]int, map[int]map[BitmapBinary]int) {
	laggingsMaps := make(map[int]map[BitmapBinary]int)
	futureMaps := make(map[int]map[BitmapBinary]int)
	lws := b.lagWindowSize
	fws := b.futureWindowSize
	chunkSize := b.chunkSize
	dimension := timeSeries.Size()

	var lwLeaveChunk, lwEnterChunk, fwLeaveChunk, fwEnterChunk BitmapBinary

	for i := 0; i < dimension; i++ {
		if (i < lws) || (i > dimension-fws) {
			laggingsMaps[i] = nil
		} else {
			if laggingsMaps[i-1] == nil {
				laggingsMaps[i] = b.constructChunkFrequencyMap(sax[i-lws : i])
				lwLeaveChunk = sax.Slice(0, chunkSize)
				lwEnterChunk = sax.Slice(i-chunkSize+1, i+1)

				futureMaps[i] = b.constructChunkFrequencyMap(sax[i : i+fws])
				fwLeaveChunk = sax.Slice(i, i+chunkSize)
				fwEnterChunk = sax.Slice(i+fws+1-chunkSize, i+fws+1)
			} else {
				lagMap := laggingsMaps[i-1]
				lagMap[lwLeaveChunk]--
				lagMap[lwEnterChunk]++
				laggingsMaps[i] = lagMap

				futureMap := futureMaps[i-1]
				futureMap[fwLeaveChunk]--
				futureMap[fwEnterChunk]++
				futureMaps[i] = futureMap

				// Update leave and enter chunks
				lwLeaveChunk = sax.Slice(i-lws, i-lws+chunkSize)
				lwEnterChunk = sax.Slice(i-chunkSize+1, i+1)
				fwLeaveChunk = sax.Slice(i, i+chunkSize)
				fwEnterChunk = sax.Slice(i+fws+1-chunkSize, i+fws+1)
			}
		}
	}
	return laggingsMaps, futureMaps
}

func (b *Bitmap) sanityCheck(timeSeries *TimeSeries) (*TimeSeries, error) {
	windowsDimension := b.lagWindowSize + b.futureWindowSize
	if (timeSeries.Size() < windowsDimension) || (windowsDimension < minimalPointsInWindows) {
		return nil, errors.New("not enough data points")
	}
	return timeSeries, nil
}
