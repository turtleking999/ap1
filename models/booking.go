package models

import (
	"time"
)

type Booking struct {
	ID             int       `json:"id"`
	PassengerID    int       `json:"passenger_id"`
	FlightID       int       `json:"flight_id"`
	Class          string    `json:"class"` // "economy", "business", "first"
	SeatNumber     string    `json:"seat_number"`
	Status         string    `json:"status"` // "confirmed", "cancelled", "checked-in"
	BookingTime    time.Time `json:"booking_time"`
	CheckInTime    time.Time `json:"check_in_time,omitempty"`
	HasCheckedIn   bool      `json:"has_checked_in"`
	Price          Money     `json:"price"`
	Compensation   Money     `json:"compensation,omitempty"`
	RiskScore      float64   `json:"risk_score"`
	IsCheapestFare bool      `json:"is_cheapest_fare"`

	// 關聯
	Passenger *Passenger `json:"passenger,omitempty"`
	Flight    *Flight    `json:"flight,omitempty"`

	// 額外信息
	SpecialRequests []string    `json:"special_requests,omitempty"`
	BaggageInfo     BaggageInfo `json:"baggage_info"`

	// 取消和退款信息
	CancellationTime time.Time `json:"cancellation_time,omitempty"`
	RefundAmount     Money     `json:"refund_amount,omitempty"`

	// 超賣相關
	IsOverbooked bool   `json:"is_overbooked"`
	UpgradedFrom string `json:"upgraded_from,omitempty"` // 如果是因超賣而升級，這裡記錄原始艙位

	// 審計字段
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type BaggageInfo struct {
	CheckedBags  int     `json:"checked_bags"`
	CarryOnBags  int     `json:"carry_on_bags"`
	TotalWeight  float64 `json:"total_weight"` // 單位：kg
	ExcessWeight float64 `json:"excess_weight,omitempty"`
	ExcessCharge Money   `json:"excess_charge,omitempty"`
}

// 用於批量操作的結構
type BookingBatch struct {
	Bookings []*Booking `json:"bookings"`
}

// 用於搜索和過濾的結構
type BookingFilter struct {
	PassengerID  int       `json:"passenger_id,omitempty"`
	FlightID     int       `json:"flight_id,omitempty"`
	Status       string    `json:"status,omitempty"`
	Class        string    `json:"class,omitempty"`
	DateFrom     time.Time `json:"date_from,omitempty"`
	DateTo       time.Time `json:"date_to,omitempty"`
	IsOverbooked *bool     `json:"is_overbooked,omitempty"`
}

// 用於更新操作的結構
type BookingUpdate struct {
	Class           *string      `json:"class,omitempty"`
	SeatNumber      *string      `json:"seat_number,omitempty"`
	Status          *string      `json:"status,omitempty"`
	SpecialRequests []string     `json:"special_requests,omitempty"`
	BaggageInfo     *BaggageInfo `json:"baggage_info,omitempty"`
}

// 在文件的其他結構體定義之後添加：

type BookingTrend struct {
	TotalBookings     int     `json:"total_bookings"`
	ConfirmedBookings int     `json:"confirmed_bookings"`
	BookingRate       float64 `json:"booking_rate"`
	AveragePrice      float64 `json:"average_price"`
	// 可以根據需求添加更多字段
}
