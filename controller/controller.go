package controller

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/dubrovin/callstats/server"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Controller - main structure
type Controller struct {
	window    *Window
	delay     time.Duration
	elements  chan int
	medians   chan int
	errors    chan error
	addDelay  chan time.Duration
	getMedian chan bool
	wg        sync.WaitGroup
	server    *server.Server
}

// NewController -
func NewController(w *Window, d time.Duration, server *server.Server) *Controller {
	return &Controller{
		window:    w,
		delay:     d,
		elements:  make(chan int, w.maxSize),
		medians:   make(chan int, w.maxSize),
		errors:    make(chan error, w.maxSize),
		addDelay:  make(chan time.Duration, w.maxSize),
		getMedian: make(chan bool, w.maxSize),
		server:    server,
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
			log.Println("readCSV: record=", record)
			for value := range record {
				i, err := strconv.Atoi(strings.Trim(record[value], "\r"))
				if err != nil {
					c.errors <- err
				} else {
					c.elements <- i
				}

			}
		case d := <-c.addDelay:
			log.Println("readCSV: addDelay=", d)
			c.delay = c.delay + d
			c.server.Delay <- c.delay
		}
	}

}

func (c *Controller) writeCsv(filePath string) {

	f, err := os.Create(filePath)
	if !os.IsExist(err) {
		os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		f, err = os.Create(filePath)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(bufio.NewWriter(f))

	// True to use \r\n as the line terminator
	w.UseCRLF = true
	defer w.Flush()

	defer c.wg.Done()
	for median := range c.medians {
		log.Println("writeCSV: median=", median)
		w.Write([]string{strconv.Itoa(median)})
		w.Flush()
	}

}

func (c *Controller) windowRunner() {
	defer c.wg.Done()
	defer close(c.medians)
	defer close(c.errors)

	for {
		select {
		case n, ok := <-c.elements:
			log.Println("windowRunner elements: ", n, ok)
			if !ok {
				return
			}
			c.window.Add(n)
			c.medians <- c.window.Median()
		case <-c.getMedian:
			c.server.Median <- c.window.Median()

		}
	}

}

func (c *Controller) errorReader() {
	defer c.wg.Done()
	for err := range c.errors {
		log.Println(err)
	}
}

// Run - runs all goroutines
func (c *Controller) Run(filePath, addr string) {
	c.wg.Add(1)
	go c.errorReader()

	c.wg.Add(1)
	go c.readCsv(filePath)

	c.wg.Add(1)
	go c.windowRunner()

	c.wg.Add(1)
	go c.writeCsv(fmt.Sprintf("out/%v_out.csv", time.Now().Unix()))
	c.wg.Add(2)
	c.registerHandlers()
	c.server.Run()
	c.wg.Wait()

}
