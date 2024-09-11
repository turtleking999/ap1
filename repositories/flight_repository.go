package repositories

import (
	"context"
	"database/sql"
	"time"

	"airline-booking/logger"
	"airline-booking/models"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type FlightRepository interface {
	SearchFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error)
}

type flightRepository struct {
	db *sql.DB
}

func NewFlightRepository(db *sql.DB) FlightRepository {
	return &flightRepository{db: db}
}

func (r *flightRepository) SearchFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlightRepository.SearchFlights")
	defer span.Finish()

	startTime := time.Now()
	defer func() {
		logger.LogWithTracing(ctx, "Flight search completed",
			zap.Duration("duration", time.Since(startTime)),
			zap.String("origin", req.Origin),
			zap.String("destination", req.Destination))
	}()

	offset := (req.Page - 1) * req.PageSize
	query := `
		SELECT id, origin, destination, departure_time, price, available_seats, total_seats
		FROM flights
		WHERE origin = $1 AND destination = $2 AND DATE(departure_time) = $3
		ORDER BY departure_time
		LIMIT $4 OFFSET $5
	`

	// 使用 context 來執行查詢
	rows, err := r.db.QueryContext(ctx, query, req.Origin, req.Destination, req.Date, req.PageSize, offset)
	if err != nil {
		logger.LogWithTracing(ctx, "Failed to execute flight search query",
			zap.Error(err),
			zap.String("origin", req.Origin),
			zap.String("destination", req.Destination))
		return nil, err
	}
	defer rows.Close()

	var flights []models.Flight
	for rows.Next() {
		var f models.Flight
		err := rows.Scan(&f.ID, &f.Origin, &f.Destination, &f.DepartureTime, &f.Price, &f.AvailableSeats, &f.TotalSeats)
		if err != nil {
			logger.LogWithTracing(ctx, "Failed to scan flight row", zap.Error(err))
			return nil, err
		}
		flights = append(flights, f)
	}

	span.SetTag("flights.count", len(flights))

	return flights, nil
}
