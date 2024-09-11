package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"airline-booking/logger"
	"airline-booking/models"
	"airline-booking/repositories"

	"go.uber.org/zap"
)

type OverbookingService interface {
	HandleOverbooking(ctx context.Context, flightID int) error
	AdjustOverbookingRatio(ctx context.Context, flightID int) error
	AssessRisk(ctx context.Context, booking *models.Booking) (float64, error)
}

type overbookingService struct {
	flightRepo    repositories.FlightRepository
	bookingRepo   repositories.BookingRepository
	notifyService NotificationService
}

func NewOverbookingService(
	flightRepo repositories.FlightRepository,
	bookingRepo repositories.BookingRepository,
	notifyService NotificationService,
) OverbookingService {
	return &overbookingService{
		flightRepo:    flightRepo,
		bookingRepo:   bookingRepo,
		notifyService: notifyService,
	}
}

func (s *overbookingService) HandleOverbooking(ctx context.Context, flightID int) error {
	flight, err := s.flightRepo.GetFlightByID(ctx, flightID)
	if err != nil {
		return err
	}

	overbooked, err := s.getOverbookedBookings(ctx, flight)
	if err != nil {
		return err
	}

	// 按風險分數排序預訂
	sort.Slice(overbooked, func(i, j int) bool {
		return overbooked[i].RiskScore > overbooked[j].RiskScore
	})

	for _, booking := range overbooked {
		if booking.Class == "economy" {
			// 嘗試升級到商務艙
			if err := s.tryUpgrade(ctx, booking, flight, "business"); err == nil {
				continue
			}
		}
		// 如果無法升級或不是經濟艙，提供補償
		if err := s.provideCompensation(ctx, booking); err != nil {
			logger.Error("Failed to provide compensation", zap.Error(err), zap.Int("bookingID", booking.ID))
		}
	}

	return nil
}

func (s *overbookingService) AdjustOverbookingRatio(ctx context.Context, flightID int) error {
	flight, err := s.flightRepo.GetFlightByID(ctx, flightID)
	if err != nil {
		return err
	}

	historicalData, err := s.flightRepo.GetHistoricalNoShowRate(ctx, flight.Route(), flight.DepartureTime.Weekday())
	if err != nil {
		return err
	}

	currentBookingTrend, err := s.bookingRepo.GetCurrentBookingTrend(ctx, flightID)
	if err != nil {
		return err
	}

	flight.EconomySeats.OverbookingRatio = s.calculateOptimalRatio(historicalData, currentBookingTrend, "economy")
	flight.BusinessSeats.OverbookingRatio = s.calculateOptimalRatio(historicalData, currentBookingTrend, "business")
	flight.FirstClassSeats.OverbookingRatio = s.calculateOptimalRatio(historicalData, currentBookingTrend, "first")

	return s.flightRepo.UpdateFlight(ctx, flight)
}

func (s *overbookingService) AssessRisk(ctx context.Context, booking *models.Booking) (float64, error) {
	riskScore := 0.0

	// 考慮多個因素來計算風險分數
	if booking.IsCheapestFare {
		riskScore += 0.3
	}

	passengerHistory, err := s.bookingRepo.GetPassengerHistory(ctx, booking.PassengerID)
	if err != nil {
		return 0, err
	}

	if passengerHistory.IsFrequentFlyer {
		riskScore -= 0.2
	}

	if booking.HasCheckedIn {
		riskScore -= 0.5
	}

	// 考慮預訂時間與起飛時間的接近程度
	timeUntilDeparture := time.Until(booking.Flight.DepartureTime)
	if timeUntilDeparture < 24*time.Hour {
		riskScore -= 0.3
	}

	return math.Max(0, math.Min(1, riskScore)), nil
}

func (s *overbookingService) getOverbookedBookings(ctx context.Context, flight *models.Flight) ([]*models.Booking, error) {
	allBookings, err := s.bookingRepo.GetBookingsByFlight(ctx, flight.ID)
	if err != nil {
		return nil, err
	}

	var overbooked []*models.Booking
	totalSeats := flight.EconomySeats.Total + flight.BusinessSeats.Total + flight.FirstClassSeats.Total

	if len(allBookings) > totalSeats {
		overbooked = allBookings[totalSeats:]
		for _, booking := range overbooked {
			riskScore, err := s.AssessRisk(ctx, booking)
			if err != nil {
				logger.Error("Failed to assess risk for booking", zap.Error(err), zap.Int("bookingID", booking.ID))
				continue
			}
			booking.RiskScore = riskScore
		}
	}

	return overbooked, nil
}

func (s *overbookingService) tryUpgrade(ctx context.Context, booking *models.Booking, flight *models.Flight, targetClass string) error {
	var availableSeats int
	switch targetClass {
	case "business":
		availableSeats = flight.BusinessSeats.Total - flight.BusinessSeats.Booked
	case "first":
		availableSeats = flight.FirstClassSeats.Total - flight.FirstClassSeats.Booked
	default:
		return errors.New("invalid target class for upgrade")
	}

	if availableSeats > 0 {
		booking.Class = targetClass
		if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
			return err
		}
		return s.notifyService.NotifyPassenger(ctx, booking, "You have been upgraded to "+targetClass+" class due to overbooking.")
	}

	return errors.New("no available seats for upgrade")
}

func (s *overbookingService) provideCompensation(ctx context.Context, booking *models.Booking) error {
	compensation := s.calculateCompensation(booking)
	booking.Compensation = compensation
	if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
		return err
	}
	return s.notifyService.NotifyPassenger(ctx, booking, fmt.Sprintf("Due to overbooking, we are offering you compensation of %v", compensation))
}

func (s *overbookingService) calculateCompensation(booking *models.Booking) models.Money {
	// 實現補償計算邏輯
	// 這可能基於航班距離、票價、乘客忠誠度等因素
	baseCompensation := booking.Price.Amount * 2 // 例如，補償為票價的兩倍
	return models.Money{Amount: baseCompensation, Currency: booking.Price.Currency}
}

func (s *overbookingService) calculateOptimalRatio(historicalData models.HistoricalData, currentTrend models.BookingTrend, class string) float64 {
	// 實現複雜的算法來計算最佳超賣比例
	// 考慮歷史 no-show 率、當前預訂趨勢、座位類型等因素
	baseRatio := historicalData.AverageNoShowRate

	switch class {
	case "economy":
		baseRatio += 0.05 // 經濟艙可以有更高的超賣比例
	case "business":
		baseRatio += 0.02 // 商務艙超賣比例稍低
	case "first":
		baseRatio += 0.01 // 頭等艙超賣比例最低
	}

	// 根據當前預訂趨勢調整比例
	if currentTrend.BookingRate > historicalData.AverageBookingRate {
		baseRatio -= 0.02 // 如果預訂率高於平均，稍微降低超賣比例
	} else {
		baseRatio += 0.02 // 如果預訂率低於平均，稍微提高超賣比例
	}

	// 確保比例在合理範圍內
	return math.Max(0, math.Min(0.2, baseRatio)) // 假設最大超賣比例為 20%
}
