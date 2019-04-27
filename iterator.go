package anomalia

import "sync"

// Iterator wraps a slice of float64 values with the current element position
type Iterator struct {
	data    []float64
	current int
	mu      sync.Mutex
}

// NewIterator returns an iterator instance
func NewIterator(data []float64) *Iterator {
	return &Iterator{data: data, current: -1}
}

// Next returns next item from the iterator
// It panics when iterator is exhausted.
func (it *Iterator) Next() *float64 {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.current++

	if it.current >= len(it.data) {
		return nil
	}
	return it.value()
}

// Value return current value of the iterator
func (it *Iterator) value() *float64 {
	return &it.data[it.current]
}
