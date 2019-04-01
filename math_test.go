package anomalia

import "testing"

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
