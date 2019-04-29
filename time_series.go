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
	length := len(sorted)
	mid := length / 2

	if length%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

// Align aligns two time series so that they have the same dimension and same timestamps
func (ts *TimeSeries) Align(other *TimeSeries) {
	var (
		it                = NewIterator(ts.Timestamps)
		otherIt           = NewIterator(other.Timestamps)
		zippedSeries      = ts.Zip()
		zippedOtherSeries = other.Zip()
		aligned           = make(map[float64]float64)
		otherAligned      = make(map[float64]float64)
	)

	timestamp, otherTimestamp := it.Next(), otherIt.Next()
	for timestamp != nil && otherTimestamp != nil {
		_timestamp, _otherTimestamp := *timestamp, *otherTimestamp
		_value, _otherValue := zippedSeries[_timestamp], zippedOtherSeries[_otherTimestamp]
		if _timestamp == _otherTimestamp {
			aligned[_timestamp] = _value
			otherAligned[_otherTimestamp] = _otherValue
			timestamp = it.Next()
			otherTimestamp = otherIt.Next()
		} else if _timestamp < _otherTimestamp {
			aligned[_timestamp] = _value
			otherAligned[_timestamp] = _otherValue
			timestamp = it.Next()
		} else {
			aligned[_otherTimestamp] = _value
			otherAligned[_otherTimestamp] = _otherValue
			otherTimestamp = otherIt.Next()
		}
	}

	//
	// Align remainder of timestamps
	//
	for timestamp != nil {
		_timestamp := *timestamp
		aligned[_timestamp] = zippedSeries[_timestamp]
		otherAligned[_timestamp] = other.Values[len(other.Values)-1]
		timestamp = it.Next()
	}

	for otherTimestamp != nil {
		_otherTimestamp := *otherTimestamp
		aligned[_otherTimestamp] = ts.Values[len(ts.Values)-1]
		otherAligned[_otherTimestamp] = zippedOtherSeries[_otherTimestamp]
		otherTimestamp = otherIt.Next()
	}

	// Adapt both the original and other time series
	alignedTimestamps, alignedValues := unpackMap(aligned)
	ts.Timestamps = alignedTimestamps
	ts.Values = alignedValues

	otherTimestamps, otherValues := unpackMap(otherAligned)
	other.Timestamps = otherTimestamps
	other.Values = otherValues
}

// Size returns the time series dimension/size.
func (ts *TimeSeries) Size() int {
	return len(ts.Timestamps)
}
