package controllers

import (
	"encoding/json"
	"time"

	"airline-booking/models"
	"airline-booking/services"

	"github.com/valyala/fasthttp"
)

type FlightController struct {
	service services.FlightService
}

func NewFlightController(service services.FlightService) *FlightController {
	return &FlightController{service: service}
}

func (c *FlightController) SearchFlights(ctx *fasthttp.RequestCtx) {
	var req models.SearchRequest
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	requestID, err := c.service.SearchFlights(ctx, req)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusAccepted)
	ctx.SetBodyString(requestID)
}

func (c *FlightController) GetSearchResults(ctx *fasthttp.RequestCtx) {
	requestID := ctx.QueryArgs().Peek("request_id")
	if requestID == nil {
		ctx.Error("Missing request_id", fasthttp.StatusBadRequest)
		return
	}

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		flights, err := c.service.GetSearchResults(string(requestID))
		if err == nil {
			ctx.SetContentType("application/json")
			json.NewEncoder(ctx).Encode(flights)
			return
		}
		time.Sleep(time.Second)
	}

	ctx.Error("Results not ready, please try again later", fasthttp.StatusNotFound)
}
