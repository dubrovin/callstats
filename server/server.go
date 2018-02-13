package server

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

// Server -
type Server struct {
	Router  *routing.Router
	errChan chan error
	Delay   chan time.Duration
	Median  chan int
	addr    string
}

// NewServer -
func NewServer(addr string, chanSize int) *Server {
	return &Server{
		Router:  routing.New(),
		errChan: make(chan error, chanSize),
		addr:    addr,
		Delay:   make(chan time.Duration, chanSize),
		Median:  make(chan int, chanSize),
	}
}

// Run -
func (c *Server) Run() {
	go c.ListenAndServe()
	go c.ReadErrChan()
}

// ListenAndServe -
func (c *Server) ListenAndServe() {
	log.Print("Listen and server addr = ", c.addr)
	c.errChan <- fasthttp.ListenAndServe(c.addr, c.Router.HandleRequest)
}

// ReadErrChan -
func (c *Server) ReadErrChan() {
	for err := range c.errChan {
		log.Print("handlers server error: ", err)
	}
}
