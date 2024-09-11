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
	GetFlightByID(ctx context.Context, flightID int) (*models.Flight, error)
	UpdateFlight(ctx context.Context, flight *models.Flight) error
	GetHistoricalNoShowRate(ctx context.Context, route string, dayOfWeek time.Weekday) (models.HistoricalData, error)
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

func (r *flightRepository) GetFlightByID(ctx context.Context, flightID int) (*models.Flight, error) {
	// 實現獲取單個航班的邏輯
	query := `
		SELECT id, origin, destination, departure_time, price, 
			   economy_seats_total, economy_seats_booked, economy_seats_overbooking_ratio,
			   business_seats_total, business_seats_booked, business_seats_overbooking_ratio,
			   first_class_seats_total, first_class_seats_booked, first_class_seats_overbooking_ratio
		FROM flights
		WHERE id = $1
	`
	var flight models.Flight
	err := r.db.QueryRowContext(ctx, query, flightID).Scan(
		&flight.ID, &flight.Origin, &flight.Destination, &flight.DepartureTime, &flight.Price,
		&flight.EconomySeats.Total, &flight.EconomySeats.Booked, &flight.EconomySeats.OverbookingRatio,
		&flight.BusinessSeats.Total, &flight.BusinessSeats.Booked, &flight.BusinessSeats.OverbookingRatio,
		&flight.FirstClassSeats.Total, &flight.FirstClassSeats.Booked, &flight.FirstClassSeats.OverbookingRatio,
	)
	if err != nil {
		return nil, err
	}
	return &flight, nil
}

func (r *flightRepository) UpdateFlight(ctx context.Context, flight *models.Flight) error {
	// 實現更新航班信息的邏輯
	query := `
		UPDATE flights
		SET origin = $2, destination = $3, departure_time = $4, price = $5,
			economy_seats_total = $6, economy_seats_booked = $7, economy_seats_overbooking_ratio = $8,
			business_seats_total = $9, business_seats_booked = $10, business_seats_overbooking_ratio = $11,
			first_class_seats_total = $12, first_class_seats_booked = $13, first_class_seats_overbooking_ratio = $14
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		flight.ID, flight.Origin, flight.Destination, flight.DepartureTime, flight.Price,
		flight.EconomySeats.Total, flight.EconomySeats.Booked, flight.EconomySeats.OverbookingRatio,
		flight.BusinessSeats.Total, flight.BusinessSeats.Booked, flight.BusinessSeats.OverbookingRatio,
		flight.FirstClassSeats.Total, flight.FirstClassSeats.Booked, flight.FirstClassSeats.OverbookingRatio,
	)
	return err
}

func (r *flightRepository) GetHistoricalNoShowRate(ctx context.Context, route string, dayOfWeek time.Weekday) (models.HistoricalData, error) {
	// 實現獲取歷史 no-show 率的邏輯
	query := `
		SELECT AVG(no_show_rate), AVG(booking_rate)
		FROM flight_statistics
		WHERE route = $1 AND day_of_week = $2
	`
	var data models.HistoricalData
	err := r.db.QueryRowContext(ctx, query, route, dayOfWeek).Scan(&data.AverageNoShowRate, &data.AverageBookingRate)
	return data, err
}
