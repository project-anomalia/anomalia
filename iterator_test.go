package anomalia

import "testing"

func TestIterator(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	it := NewIterator(data)
	pos := 0
	for {
		if valuePtr := it.Next(); valuePtr != nil {
			if data[pos] != *valuePtr {
				t.Fatal("iterator value mismatch")
			}
			pos++
			continue
		}
		break
	}
}
