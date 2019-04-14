package anomalia

import "unicode/utf8"

// Binary wrapper type around a string with custom behaviour
type BitmapBinary string

// Slice slices a string in a pythonic way
// When lower == upper, it returns an empty string
// When lower < 0 and upper < len(binary), it return an empty string
// When lower < 0 and upper >= len(binary), it returns the first character
// When lower >= 0 and upper >= len(binary), it slices the string from lower till end of string
func (bb BitmapBinary) Slice(lower, upper int) BitmapBinary {
	var result string
	switch {
	case lower == upper:
	case (lower < 0) && (upper < len(bb)):
		result = ""
	case (lower < 0) && (upper >= len(bb)):
		result = bb.String()[0:1]
	case (lower >= 0) && (upper >= len(bb)):
		result = bb.String()[lower:]
	default:
		result = bb.String()[lower:upper]
	}
	return BitmapBinary(result)
}

// At returns character string at the specified index
func (bb BitmapBinary) At(index int) BitmapBinary {
	if index >= len(bb) {
		return ""
	}
	return BitmapBinary(bb.String()[index])
}

// String returns the underlaying string
func (bb BitmapBinary) String() string {
	return string(bb)
}

// Len returns the length of the underlaying string
func (bb BitmapBinary) Len() int {
	return utf8.RuneCountInString(bb.String())
}
