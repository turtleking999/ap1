package main

import (
	"airline-booking/config"
	"airline-booking/controllers"
	"airline-booking/logger"
	"airline-booking/repositories"
	"airline-booking/routes"
	"airline-booking/services"

	"github.com/fasthttp/router"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Init(); err != nil {
		panic(err)
	}

	if err := logger.InitTracer("airline-booking"); err != nil {
		logger.Fatal("Failed to initialize tracer", zap.Error(err))
	}

	cfg := config.NewConfig()

	db, err := config.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	redisClient, err := config.InitRedis(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize Redis", zap.Error(err))
	}
	defer redisClient.Close()

	flightRepo := repositories.NewFlightRepository(db)
	flightService := services.NewFlightService(flightRepo, redisClient)
	flightController := controllers.NewFlightController(flightService)

	go flightService.ProcessSearchRequests()

	r := router.New()
	routes.SetupRoutes(r, flightController)

	handler := func(ctx *fasthttp.RequestCtx) {
		span, traceCtx := opentracing.StartSpanFromContext(ctx, "http_handler")
		defer span.Finish()

		r.Handler(ctx)

		logger.LogWithTracing(traceCtx, "Request processed",
			zap.String("method", string(ctx.Method())),
			zap.String("path", string(ctx.Path())),
			zap.Int("status", ctx.Response.StatusCode()))
	}

	logger.Info("Server is running", zap.String("port", cfg.ServerPort))
	if err := fasthttp.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		logger.Fatal("Server stopped", zap.Error(err))
	}
}
