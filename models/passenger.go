package models

import (
	"time"
)

type Passenger struct {
	ID             int       `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phone_number"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	Nationality    string    `json:"nationality"`
	PassportNumber string    `json:"passport_number"`
	PassportExpiry time.Time `json:"passport_expiry"`

	// 地址信息
	Address    string `json:"address"`
	City       string `json:"city"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`

	// 忠誠度計劃信息
	FrequentFlyerNumber string `json:"frequent_flyer_number,omitempty"`
	FrequentFlyerTier   string `json:"frequent_flyer_tier,omitempty"` // 例如：Silver, Gold, Platinum
	FrequentFlyerPoints int    `json:"frequent_flyer_points,omitempty"`

	// 特殊需求和偏好
	SpecialMealPreference string `json:"special_meal_preference,omitempty"`
	SeatPreference        string `json:"seat_preference,omitempty"` // 例如：Window, Aisle, Extra Legroom
	SpecialAssistance     bool   `json:"special_assistance"`

	// 統計和分析
	TotalFlights   int       `json:"total_flights"`
	TotalSpent     float64   `json:"total_spent"`
	LastFlightDate time.Time `json:"last_flight_date,omitempty"`

	// 行銷相關
	MarketingConsent  bool   `json:"marketing_consent"`
	PreferredLanguage string `json:"preferred_language"`

	// 審計字段
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 用於搜索和過濾的結構
type PassengerFilter struct {
	FirstName           string    `json:"first_name,omitempty"`
	LastName            string    `json:"last_name,omitempty"`
	Email               string    `json:"email,omitempty"`
	FrequentFlyerNumber string    `json:"frequent_flyer_number,omitempty"`
	Nationality         string    `json:"nationality,omitempty"`
	FrequentFlyerTier   string    `json:"frequent_flyer_tier,omitempty"`
	MinTotalFlights     int       `json:"min_total_flights,omitempty"`
	MinTotalSpent       float64   `json:"min_total_spent,omitempty"`
	LastFlightAfter     time.Time `json:"last_flight_after,omitempty"`
}

// 用於更新操作的結構
type PassengerUpdate struct {
	Email                 *string    `json:"email,omitempty"`
	PhoneNumber           *string    `json:"phone_number,omitempty"`
	Address               *string    `json:"address,omitempty"`
	City                  *string    `json:"city,omitempty"`
	Country               *string    `json:"country,omitempty"`
	PostalCode            *string    `json:"postal_code,omitempty"`
	PassportNumber        *string    `json:"passport_number,omitempty"`
	PassportExpiry        *time.Time `json:"passport_expiry,omitempty"`
	FrequentFlyerTier     *string    `json:"frequent_flyer_tier,omitempty"`
	SpecialMealPreference *string    `json:"special_meal_preference,omitempty"`
	SeatPreference        *string    `json:"seat_preference,omitempty"`
	SpecialAssistance     *bool      `json:"special_assistance,omitempty"`
	MarketingConsent      *bool      `json:"marketing_consent,omitempty"`
	PreferredLanguage     *string    `json:"preferred_language,omitempty"`
}

// 在現有的 Passenger 結構體之後添加：

type PassengerHistory struct {
	PassengerID              int                 `json:"passenger_id"`
	IsFrequentFlyer          bool                `json:"is_frequent_flyer"`
	TotalFlights             int                 `json:"total_flights"`
	TotalSpent               float64             `json:"total_spent"`
	LastFlightDate           time.Time           `json:"last_flight_date"`
	AverageFlightFrequency   float64             `json:"average_flight_frequency"` // 例如：每月平均飛行次數
	MostFrequentRoute        string              `json:"most_frequent_route"`
	PreferredAirlinePartners []string            `json:"preferred_airline_partners"`
	LoyaltyTierHistory       []LoyaltyTierChange `json:"loyalty_tier_history"`
	SpecialServiceRequests   []string            `json:"special_service_requests"` // 歷史上的特殊服務請求
	CancellationRate         float64             `json:"cancellation_rate"`        // 取消預訂的比率
	OnTimeCheckInRate        float64             `json:"on_time_check_in_rate"`    // 準時登機的比率
	FeedbackScore            float64             `json:"feedback_score"`           // 平均客戶反饋分數
	LastUpdated              time.Time           `json:"last_updated"`
}

type LoyaltyTierChange struct {
	Tier       string    `json:"tier"`
	ChangeDate time.Time `json:"change_date"`
}
