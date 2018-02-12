package controller

import (
	"sort"
	"time"
)

// Window - sliding window struct
type Window struct {
	nodes   map[int64]int
	maxSize int
}

// NewWindow -
func NewWindow(maxSize int) *Window {
	return &Window{
		maxSize: maxSize,
		nodes:   make(map[int64]int, maxSize),
	}
}

// minKey - finds min key in nodes
func (w *Window) minKey() (min int64) {
	min = time.Now().UnixNano()
	for k := range w.nodes {
		if k < min {
			min = k
		}
	}
	return min
}

// Add -
func (w *Window) Add(node int) {
	if len(w.nodes)+1 > w.maxSize {
		delete(w.nodes, w.minKey())
	}
	w.nodes[time.Now().UnixNano()] = node

}

// Median -
func (w *Window) Median() (median int) {
	n := len(w.nodes)
	if n == 1 {
		return -1
	}

	var durations []int
	for _, v := range w.nodes {
		durations = append(durations, v)
	}
	sort.Ints(durations)

	if n%2 == 0 {
		// even
		median = (durations[n/2-1] + durations[n/2]) / 2
	} else {
		// odd
		median = durations[(n)/2]
	}

	return median

}
