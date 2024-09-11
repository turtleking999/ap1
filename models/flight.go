package models

import "time"

type Flight struct {
	ID             int       `json:"id"`
	Origin         string    `json:"origin"`
	Destination    string    `json:"destination"`
	DepartureTime  time.Time `json:"departure_time"`
	Price          float64   `json:"price"`
	AvailableSeats int       `json:"available_seats"`
	TotalSeats     int       `json:"total_seats"`
}

type SearchRequest struct {
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Date        time.Time `json:"date"`
	Page        int       `json:"page"`
	PageSize    int       `json:"page_size"`
}
