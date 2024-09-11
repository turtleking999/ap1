package main

import (
	"log"

	"airline-booking/config"
	"airline-booking/controllers"
	"airline-booking/repositories"
	"airline-booking/services"

	"github.com/valyala/fasthttp"
)

func main() {

	cfg := config.NewConfig()

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient, err := config.InitRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	flightRepo := repositories.NewFlightRepository(db)

	flightService := services.NewFlightService(flightRepo, redisClient)

	flightController := controllers.NewFlightController(flightService)

	handler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/search":
			flightController.SearchFlights(ctx)
		default:
			ctx.Error("Not found", fasthttp.StatusNotFound)
		}
	}

	log.Printf("Server is running on port %s", cfg.ServerPort)
	log.Fatal(fasthttp.ListenAndServe(":"+cfg.ServerPort, handler))
}
