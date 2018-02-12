package controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test1(t *testing.T) {
	// Nodes:  100, 102, 101, 110, 120, 115
	// Medians: -1, 101, 101, 102, 110, 115
	maxSize := 3
	w := NewWindow(maxSize)
	w.Add(100)
	assert.Equal(t, -1, w.Median())
	w.Add(102)
	assert.Equal(t, 101, w.Median())
	w.Add(101)
	assert.Equal(t, 101, w.Median())
	w.Add(110)
	assert.Equal(t, 102, w.Median())
	w.Add(120)
	assert.Equal(t, 110, w.Median())
	w.Add(115)
	assert.Equal(t, 115, w.Median())
}
