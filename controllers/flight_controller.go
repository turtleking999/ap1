package controllers

import (
	"encoding/json"

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

	flights, err := c.service.SearchFlights(req)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(flights)
}
