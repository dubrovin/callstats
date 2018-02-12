package controller

import (
	"testing"
	"time"
)

//
//func TestNewController(t *testing.T) {
//	maxSize := 3
//	delay := time.Second
//	w := NewWindow(maxSize)
//
//	c := NewController(w, delay)
//	c.Run("test1.csv")
//}

func TestRead(t *testing.T) {
	maxSize := 3
	delay := time.Second
	w := NewWindow(maxSize)

	c := NewController(w, delay)
	c.Run("test1.csv")
	//time.Sleep(time.Second * 15)
}
