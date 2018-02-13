package main

import (
	"flag"
	"github.com/dubrovin/callstats/controller"
	"github.com/dubrovin/callstats/server"
	"log"
	"time"
)

var (
	addr       = flag.String("addr", ":8080", "http service address")
	delay      = flag.String("delay", "1µs", "interval delay")
	windowSize = flag.Int("window_size", 100, "size for sliding window")
	filePath   = flag.String("file_path", "tests/test2.csv", "path to test file")
)

func main() {
	flag.Parse()
	intervalDelay, err := time.ParseDuration(*delay)
	if err != nil {
		log.Fatal(err)
	}
	maxSize := *windowSize
	w := controller.NewWindow(maxSize)
	s := server.NewServer(*addr, maxSize)
	c := controller.NewController(w, intervalDelay, s)
	c.Run(*filePath, *addr)

}
