package anomalia

// TimeSeries wrapper for timestamps and their values
type TimeSeries struct {
	Timestamps []float64
	Values     []float64
}

// NewTimeSeries creates a new time series data structure
func NewTimeSeries(timestamps []float64, values []float64) *TimeSeries {
	if len(timestamps) != len(values) {
		panic("timestamps and values must have the same size")
	}
	return &TimeSeries{
		Timestamps: timestamps,
		Values:     values,
	}
}

// EarliestTimestamp returns the earliest timestamp in the time series
func (ts *TimeSeries) EarliestTimestamp() float64 {
	min, _ := minMax(ts.Timestamps)
	return min
}

// LastestTimestamp returns the latest timestamp in the time series
func (ts *TimeSeries) LastestTimestamp() float64 {
	_, max := minMax(ts.Timestamps)
	return max
}

// Zip convert the time series to a map (map[Timestamp]Value)
func (ts *TimeSeries) Zip() map[float64]float64 {
	m := make(map[float64]float64)
	sorted := sortedCopy(ts.Timestamps)

	for idx, timestamp := range sorted {
		m[timestamp] = ts.Values[idx]
	}
	return m
}

// AddOffset increments time series timestamps by some offset
func (ts *TimeSeries) AddOffset(offset float64) *TimeSeries {
	offsettedTimestamps := mapSlice(ts.Timestamps, func(timestamp float64) float64 { return timestamp + offset })
	return NewTimeSeries(offsettedTimestamps, ts.Values)
}

// Normalize normalizes the time series values by dividing by the maximum value
func (ts *TimeSeries) Normalize() *TimeSeries {
	_, max := minMax(ts.Values)
	normalizedValues := mapSlice(ts.Values, func(value float64) float64 { return value / max })
	return NewTimeSeries(ts.Timestamps, normalizedValues)
}

// NormalizeWithMinMax normalizes time series values using MixMax
func (ts *TimeSeries) NormalizeWithMinMax() *TimeSeries {
	normalizedValues := ts.Values
	if min, max := minMax(ts.Values); min != max {
		normalizedValues = mapSlice(ts.Values, func(value float64) float64 { return value - min/max - min })
	}
	return NewTimeSeries(ts.Timestamps, normalizedValues)
}

// Crop crops the time series timestamps into the specified range [start, end]
func (ts *TimeSeries) Crop(start, end float64) *TimeSeries {
	zippedSeries := ts.Zip()
	// Filter timestamps within the crop range
	timestamps := filter(ts.Timestamps, func(timestamp float64) bool {
		return (timestamp >= start) && (timestamp <= end)
	})

	// Get values of cropped timestamps
	values := make([]float64, 0, len(timestamps))
	for _, timestamp := range timestamps {
		values = append(values, zippedSeries[timestamp])
	}
	return NewTimeSeries(timestamps, values)
}

// Average calculates average value over the time series
func (ts *TimeSeries) Average() float64 {
	return Average(ts.Values)
}

// Stdev calculates the standard deviation of the time series
func (ts *TimeSeries) Stdev() float64 {
	return Stdev(ts.Values)
}

// Median calculates median value over the time series.
func (ts *TimeSeries) Median() float64 {
	sorted := sortedCopy(ts.Values)
	len := len(sorted)
	mid := len / 2

	if len%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}
