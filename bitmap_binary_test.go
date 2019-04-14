package anomalia

import "testing"

func TestBitmapBinarySlice(t *testing.T) {
	var str BitmapBinary = "hello"
	slice := str.Slice(-1, 10)
	if slice != "h" {
		t.Fatalf("must return an empty string")
	}
}

func TestBitmapBinaryAtIndex(t *testing.T) {
	var str BitmapBinary = "binary"
	s := str.At(0)
	if s != "b" {
		t.Fatalf("must return the first character of the binary")
	}

	s = str.At(10)
	if s != "" {
		t.Fatalf("must return an empty string when index over binary length")
	}
}
