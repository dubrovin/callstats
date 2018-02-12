package controller

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"log"
)

// Controller - main structure
type Controller struct {
	window   *Window
	delay    time.Duration
	elements chan int
	medians  chan int
	errors   chan error
	addDelay chan time.Duration
	wg       sync.WaitGroup
}

// NewController -
func NewController(w *Window, d time.Duration) *Controller {
	return &Controller{
		window:   w,
		delay:    d,
		elements: make(chan int, w.maxSize),
		medians:  make(chan int, w.maxSize),
		errors:   make(chan error, w.maxSize),
		addDelay: make(chan time.Duration),
	}
}

func (c *Controller) readCsv(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	fmt.Println("RUN READ")
	defer c.wg.Done()

	ticker := time.NewTicker(c.delay)
	for {
		select {
		case <-ticker.C:
			record, err := r.Read()
			// Stop at EOF.
			if err == io.EOF {
				close(c.elements)
				return
			}
			for value := range record {
				i, err := strconv.Atoi(strings.Trim(record[value], "\r"))
				if err != nil {
					c.errors <- err
				} else {
					c.elements <- i
				}

			}
		case d := <-c.addDelay:
			c.delay = c.delay + d
		}
	}

}

func (c *Controller) writeCsv(filePath string) {
	f, _ := os.Create(filePath)
	defer f.Close()
	w := csv.NewWriter(bufio.NewWriter(f))

	// True to use \r\n as the line terminator
	w.UseCRLF = true
	defer w.Flush()

	defer c.wg.Done()
	for median := range c.medians {
		w.Write([]string{fmt.Sprintf("%s%s", strconv.Itoa(median), "")})
	}

}

func (c *Controller) windowRunner() {
	defer c.wg.Done()
	defer close(c.medians)
	defer close(c.errors)
	for n := range c.elements {
		c.window.Add(n)
		c.medians <- c.window.Median()
	}
}

func (c *Controller) errorReader() {
	fmt.Println("RUN ERROR READER")
	defer c.wg.Done()
	for err := range c.errors {
		fmt.Println(err)
	}
}

// Run - runs all gorutines
func (c *Controller) Run(filePath string) {
	c.wg.Add(1)
	go c.errorReader()

	c.wg.Add(1)
	go c.readCsv(filePath)

	c.wg.Add(1)
	go c.windowRunner()

	c.wg.Add(1)
	go c.writeCsv(fmt.Sprintf("tests/%v_out.csv", time.Now().Unix()))

	c.wg.Wait()
}
