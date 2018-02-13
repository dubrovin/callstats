package controller

import (
	"encoding/json"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"time"
)

// AddDelayRequest -
type AddDelayRequest struct {
	Delay string `json:"delay"`
}

// AddDelayResponse -
type AddDelayResponse struct {
	Delay time.Duration `json:"currentDelay"`
}

// GetMedianResp -
type GetMedianResp struct {
	Median int `json:"median"`
}

func (c *Controller) registerHandlers() {
	c.server.Router.Get("/get_median", c.GetMedian)
	c.server.Router.Post("/add_delay", c.AddDelay)
}

// GetMedian -
func (c *Controller) GetMedian(ctx *routing.Context) error {
	ctx.SetContentType("application/json")

	c.getMedian <- true
	median := <-c.server.Median

	jsonData, err := json.Marshal(GetMedianResp{Median: median})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	err = ctx.WriteData(jsonData)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

// AddDelay -
func (c *Controller) AddDelay(ctx *routing.Context) error {
	ctx.SetContentType("application/json")

	requestData := &AddDelayRequest{}
	err := json.Unmarshal(ctx.PostBody(), requestData)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	delay, err := time.ParseDuration(requestData.Delay)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}
	c.addDelay <- delay
	currentDelay := <-c.server.Delay

	jsonData, err := json.Marshal(AddDelayResponse{Delay: currentDelay})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	err = ctx.WriteData(jsonData)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}
