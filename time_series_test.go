package anomalia

import (
	"reflect"
	"testing"
)

var (
	timestamps = []float64{2, 3, 5, 8, 9, 10, 15}
	values     = []float64{1.0, 0.6, 2.5, 1.9, 0.3, 3.2, 0}
)

func TestNewTimeSeries(t *testing.T) {
	ts := NewTimeSeries(timestamps, values)
	actualType := "*anomalia.TimeSeries"
	expectedType := reflect.TypeOf(ts).String()
	if expectedType != actualType {
		t.Fatalf("expected '%s', got '%s'", expectedType, actualType)
	}

	// Assert panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewTimeSeries did not panic")
		}
	}()
	NewTimeSeries([]float64{1, 2}, []float64{1})
}

func TestEarliestTimestamp(t *testing.T) {
	ts := NewTimeSeries(timestamps, values)
	actual := ts.EarliestTimestamp()
	expected := 2.0
	if actual != expected {
		t.Fatalf("expected '%v', got '%v'", expected, actual)
	}
}

func TestLastestTimestamp(t *testing.T) {
	ts := NewTimeSeries(timestamps, values)
	actual := ts.LastestTimestamp()
	expected := 15.0
	if actual != expected {
		t.Fatalf("expected '%v', got '%v'", expected, actual)
	}
}

func TestZipTimeSeries(t *testing.T) {
	ts := NewTimeSeries(timestamps, values)
	expected := ts.Zip()
	if (len(expected) != len(ts.Timestamps)) || (len(expected) != len(ts.Values)) {
		t.Fatalf("time series lengths do not match")
	}
}

func TestAddOffsetToTimeSeries(t *testing.T) {
	ts := NewTimeSeries(timestamps, values)
	offsetted := ts.AddOffset(1)
	if len(ts.Timestamps) != len(offsetted.Timestamps) {
		t.Fatalf("offsetted time series length do not match")
	}
}

func TestNormalize(t *testing.T) {
	ts := NewTimeSeries(timestamps, values).Normalize()
	if ts == nil {
		t.Fatalf("normalized time series cannot be nil")
	}
}

func TestNormalizeWithMinMax(t *testing.T) {
	ts := NewTimeSeries(timestamps, values).NormalizeWithMinMax()
	if ts == nil {
		t.Fatalf("minMax normalized time series cannot be nil")
	}
}

func TestCrop(t *testing.T) {
	ts := NewTimeSeries(timestamps, values).Crop(0, 4)
	if len(ts.Timestamps) != 2 || len(ts.Values) != 2 {
		t.Fatalf("expected size to be 2, got %v", len(ts.Timestamps))
	}
}

func TestTimeSeriesAverage(t *testing.T) {
	actual := Float64WithPrecision(NewTimeSeries(timestamps, values).Average(), 2)
	expected := Float64WithPrecision(1.36, 2)
	if actual != expected {
		t.Fatalf("expected %f, got %f", expected, actual)
	}
}

func TestMedian(t *testing.T) {
	ts := NewTimeSeries(timestamps, values)

	actual := Float64WithPrecision(ts.Median(), 2)
	expected := Float64WithPrecision(1.00, 2)
	if expected != actual {
		t.Fatalf("expected %f, got %f", expected, actual)
	}

	ts = ts.Crop(0, 8)
	actual = Float64WithPrecision(ts.Median(), 2)
	expected = Float64WithPrecision(1.45, 2)
	if expected != actual {
		t.Fatalf("expected %f, got %f", expected, actual)
	}
}

func TestAlign1(t *testing.T) {
	ts := NewTimeSeries([]float64{4, 5, 6, 7, 8, 15}, []float64{1.2, 0, 1, 0.5, 4, 7})
	otherTs := NewTimeSeries([]float64{1, 2, 3}, []float64{0.9, 10.1, 5.4})

	ts.Align(otherTs)

	if ts.Size() != otherTs.Size() {
		t.Fatalf("time series size mismatch")
	}
}

func TestAlign2(t *testing.T) {
	ts := NewTimeSeries([]float64{1, 2, 3, 4}, []float64{0.1, 0.2, 0.3, 9.8})
	otherTs := NewTimeSeries([]float64{4, 5, 6, 7, 8, 15}, []float64{1.2, 0, 1, 0.5, 4, 7})

	ts.Align(otherTs)

	if ts.Size() != otherTs.Size() {
		t.Fatalf("time series size mismatch")
	}
}
