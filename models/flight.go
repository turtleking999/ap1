package models

import (
	"fmt"
	"time"
)

type Flight struct {
	ID             int       `json:"id"`
	Origin         string    `json:"origin"`
	Destination    string    `json:"destination"`
	DepartureTime  time.Time `json:"departure_time"`
	Price          float64   `json:"price"`
	AvailableSeats int       `json:"available_seats"`
	TotalSeats     int       `json:"total_seats"`
	EconomySeats   struct {
		Total            int
		Booked           int
		OverbookingRatio float64
	}
	BusinessSeats struct {
		Total            int
		Booked           int
		OverbookingRatio float64
	}
	FirstClassSeats struct {
		Total            int
		Booked           int
		OverbookingRatio float64
	}
}

type SearchRequest struct {
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Date        time.Time `json:"date"`
	Page        int       `json:"page"`
	PageSize    int       `json:"page_size"`
}

// 在現有的 Flight 結構體之後添加：

type HistoricalData struct {
	Route              string       `json:"route"`
	DayOfWeek          time.Weekday `json:"day_of_week"`
	AverageNoShowRate  float64      `json:"average_no_show_rate"`
	AverageBookingRate float64      `json:"average_booking_rate"`
	AverageLoadFactor  float64      `json:"average_load_factor"`
	AverageYield       float64      `json:"average_yield"`
	PeakSeasonStart    time.Time    `json:"peak_season_start,omitempty"`
	PeakSeasonEnd      time.Time    `json:"peak_season_end,omitempty"`
	LastUpdated        time.Time    `json:"last_updated"`
}

// Route 返回航班的航線字符串
func (f *Flight) Route() string {
	return fmt.Sprintf("%s-%s", f.Origin, f.Destination)
}
