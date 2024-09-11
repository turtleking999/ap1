package repositories

import (
	"database/sql"

	"airline-booking/models"
)

type FlightRepository interface {
	SearchFlights(req models.SearchRequest) ([]models.Flight, error)
}

type flightRepository struct {
	db *sql.DB
}

func NewFlightRepository(db *sql.DB) FlightRepository {
	return &flightRepository{db: db}
}

func (r *flightRepository) SearchFlights(req models.SearchRequest) ([]models.Flight, error) {
	offset := (req.Page - 1) * req.PageSize
	query := `
		SELECT id, origin, destination, departure_time, price, available_seats, total_seats
		FROM flights
		WHERE origin = $1 AND destination = $2 AND DATE(departure_time) = $3
		ORDER BY departure_time
		LIMIT $4 OFFSET $5
	`
	rows, err := r.db.Query(query, req.Origin, req.Destination, req.Date, req.PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flights []models.Flight
	for rows.Next() {
		var f models.Flight
		err := rows.Scan(&f.ID, &f.Origin, &f.Destination, &f.DepartureTime, &f.Price, &f.AvailableSeats, &f.TotalSeats)
		if err != nil {
			return nil, err
		}
		flights = append(flights, f)
	}

	return flights, nil
}
