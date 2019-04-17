package anomalia

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const minimalPointsInWindows = 50

type Bitmap struct {
	ChunkSize        int
	Precision        int
	LagWindowSize    int
	FutureWindowSize int
}

func (bit *Bitmap) Run(timeSeries *TimeSeries) *ScoreList {
	scoreList, _ := bit.computeScores(timeSeries)
	return scoreList
}

func (bit *Bitmap) computeScores(timeSeries *TimeSeries) (*ScoreList, error) {
	if _, err := bit.sanityCheck(timeSeries); err != nil {
		return nil, err
	}
	sax := bit.generateSAX(timeSeries)
	laggingsMaps, futureMaps := bit.constructAllSAXChunks(timeSeries, sax)
	dimension := len(timeSeries.Timestamps)

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
		if (idx < bit.LagWindowSize) || (idx > dimension-bit.FutureWindowSize) {
			return 0.0
		}
		return computeScoreBetweenTwoWindows(idx)
	})

	scoreList := &ScoreList{timeSeries.Timestamps, scores}
	return scoreList, nil
}

// generateSAX generates the SAX representation of the time series values
func (bit *Bitmap) generateSAX(timeSeries *TimeSeries) BitmapBinary {
	sections := make(map[int]float64)
	min, max := minMax(timeSeries.Values)

	// Break the whole value range into different sections
	sectionHeight := (max - min) / float64(bit.Precision)
	for i := 0; i < bit.Precision; i++ {
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
		return fmt.Sprintf("%v", sax)
	}

	var saxBuilder strings.Builder
	for _, value := range timeSeries.Values {
		singleSAX := generateSingleSAX(value)
		saxBuilder.WriteString(singleSAX)
	}
	return BitmapBinary(saxBuilder.String())
}

func (bit *Bitmap) constructChunkFrequencyMap(sax BitmapBinary) map[BitmapBinary]int {
	frequencyMap := make(map[BitmapBinary]int)
	saxLength := sax.Len()
	for i := 0; i < saxLength; i++ {
		if i+bit.ChunkSize <= saxLength {
			chunk := sax.Slice(i, i+bit.ChunkSize)
			frequencyMap[chunk] += 1
		}
	}
	return frequencyMap
}

func (bit *Bitmap) constructAllSAXChunks(timeSeries *TimeSeries, sax BitmapBinary) (map[int]map[BitmapBinary]int, map[int]map[BitmapBinary]int) {
	laggingsMaps := make(map[int]map[BitmapBinary]int)
	futureMaps := make(map[int]map[BitmapBinary]int)
	lws := bit.LagWindowSize
	fws := bit.FutureWindowSize
	chunkSize := bit.ChunkSize
	dimension := len(timeSeries.Values)

	var lwLeaveChunk, lwEnterChunk, fwLeaveChunk, fwEnterChunk BitmapBinary

	for i := 0; i < dimension; i++ {
		if (i < lws) || (i > dimension-fws) {
			laggingsMaps[i] = nil
		} else {
			if laggingsMaps[i-1] == nil {
				laggingsMaps[i] = bit.constructChunkFrequencyMap(sax[i-lws : i])
				lwLeaveChunk = sax.Slice(0, chunkSize)
				lwEnterChunk = sax.Slice(i-chunkSize+1, i+1)

				futureMaps[i] = bit.constructChunkFrequencyMap(sax[i : i+fws])
				fwLeaveChunk = sax.Slice(i, i+chunkSize)
				fwEnterChunk = sax.Slice(i+fws+1-chunkSize, i+fws+1)
			} else {
				lagMap := laggingsMaps[i-1]
				lagMap[lwLeaveChunk] -= 1
				lagMap[lwEnterChunk] += 1
				laggingsMaps[i] = lagMap

				futureMap := futureMaps[i-1]
				futureMap[fwLeaveChunk] -= 1
				futureMap[fwEnterChunk] += 1
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

func (bit *Bitmap) sanityCheck(timeSeries *TimeSeries) (*TimeSeries, error) {
	windowsDimension := bit.LagWindowSize + bit.FutureWindowSize
	if (len(timeSeries.Timestamps) < windowsDimension) || (windowsDimension < minimalPointsInWindows) {
		return nil, errors.New("not enough data points")
	}
	return timeSeries, nil
}
