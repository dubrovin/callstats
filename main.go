package main

import (
	"github.com/dubrovin/callstats/controller"
	"time"
)

var ()

func main() {

	maxSize := 10000
	delay := time.Nanosecond
	w := controller.NewWindow(maxSize)

	c := controller.NewController(w, delay)
	c.Run("tests/test4.csv")

}
