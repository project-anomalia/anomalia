package anomalia

import (
	"math"
	"testing"
)

var input = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func TestAverage(t *testing.T) {
	actual := Average(input)
	expected := 5.5
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestSum(t *testing.T) {
	actual := Sum(input)
	expected := 55.0
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestVariance(t *testing.T) {
	actual := Variance(input)
	expected := 8.25
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestStdev(t *testing.T) {
	actual := Stdev(input)
	expected := math.Sqrt(8.25)
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestRoundFloat(t *testing.T) {
	actual := RoundFloat(0.5)
	expected := 1
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestFloat64WithPrecision(t *testing.T) {
	actual := Float64WithPrecision(1.45424, 2)
	expected := 1.45
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestPdf(t *testing.T) {
	actual := Pdf(0.0, 1.0)(1.0)
	expected := 0.24197072451914337
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestCdf(t *testing.T) {
	actual := Cdf(0.0, 1.0)(1.0)
	expected := 0.8413447361676363
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestErf(t *testing.T) {
	actual := Erf(1.0)
	expected := 0.8427006897475899
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}

	actual = Erf(-1.0)
	expected = -expected
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
