package services

import (
	"context"
	"errors"
	"time"

	"airline-booking/logger"
	"airline-booking/models"
	"airline-booking/repositories"

	"go.uber.org/zap"
)

type BookingService interface {
	CreateBooking(ctx context.Context, booking *models.Booking) error
	GetBooking(ctx context.Context, bookingID int) (*models.Booking, error)
	UpdateBooking(ctx context.Context, booking *models.Booking) error
	CancelBooking(ctx context.Context, bookingID int) error
	ListBookingsByPassenger(ctx context.Context, passengerID int) ([]*models.Booking, error)
	CheckIn(ctx context.Context, bookingID int) error
}

type bookingService struct {
	bookingRepo        repositories.BookingRepository
	flightRepo         repositories.FlightRepository
	passengerRepo      repositories.PassengerRepository
	overbookingService OverbookingService
	notifyService      NotificationService
}

func NewBookingService(
	bookingRepo repositories.BookingRepository,
	flightRepo repositories.FlightRepository,
	passengerRepo repositories.PassengerRepository,
	overbookingService OverbookingService,
	notifyService NotificationService,
) BookingService {
	return &bookingService{
		bookingRepo:        bookingRepo,
		flightRepo:         flightRepo,
		passengerRepo:      passengerRepo,
		overbookingService: overbookingService,
		notifyService:      notifyService,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, booking *models.Booking) error {
	// 檢查航班是否存在並有足夠的座位
	flight, err := s.flightRepo.GetFlightByID(ctx, booking.FlightID)
	if err != nil {
		return err
	}

	availableSeats := s.calculateAvailableSeats(flight, booking.Class)
	if availableSeats <= 0 {
		return errors.New("no available seats")
	}

	// 檢查乘客是否存在
	_, err = s.passengerRepo.GetPassengerByID(ctx, booking.PassengerID)
	if err != nil {
		return err
	}

	// 創建預訂
	booking.Status = "confirmed"
	booking.BookingTime = time.Now()

	err = s.bookingRepo.CreateBooking(ctx, booking)
	if err != nil {
		return err
	}

	// 更新航班座位信息
	switch booking.Class {
	case "economy":
		flight.EconomySeats.Booked++
	case "business":
		flight.BusinessSeats.Booked++
	case "first":
		flight.FirstClassSeats.Booked++
	}
	err = s.flightRepo.UpdateFlight(ctx, flight)
	if err != nil {
		// 如果更新失敗，需要回滾預訂
		s.bookingRepo.DeleteBooking(ctx, booking.ID)
		return err
	}

	// 評估風險並設置風險分數
	riskScore, err := s.overbookingService.AssessRisk(ctx, booking)
	if err != nil {
		logger.Error("Failed to assess booking risk", zap.Error(err), zap.Int("bookingID", booking.ID))
	} else {
		booking.RiskScore = riskScore
		s.bookingRepo.UpdateBooking(ctx, booking)
	}

	// 發送預訂確認通知
	s.notifyService.NotifyPassenger(ctx, booking, "Your booking has been confirmed.")

	return nil
}

func (s *bookingService) GetBooking(ctx context.Context, bookingID int) (*models.Booking, error) {
	return s.bookingRepo.GetBookingByID(ctx, bookingID)
}

func (s *bookingService) UpdateBooking(ctx context.Context, booking *models.Booking) error {
	existingBooking, err := s.bookingRepo.GetBookingByID(ctx, booking.ID)
	if err != nil {
		return err
	}

	// 檢查是否需要更改座位類型
	if existingBooking.Class != booking.Class {
		flight, err := s.flightRepo.GetFlightByID(ctx, booking.FlightID)
		if err != nil {
			return err
		}

		availableSeats := s.calculateAvailableSeats(flight, booking.Class)
		if availableSeats <= 0 {
			return errors.New("no available seats in the new class")
		}

		// 更新航班座位信息
		switch existingBooking.Class {
		case "economy":
			flight.EconomySeats.Booked--
		case "business":
			flight.BusinessSeats.Booked--
		case "first":
			flight.FirstClassSeats.Booked--
		}

		switch booking.Class {
		case "economy":
			flight.EconomySeats.Booked++
		case "business":
			flight.BusinessSeats.Booked++
		case "first":
			flight.FirstClassSeats.Booked++
		}

		err = s.flightRepo.UpdateFlight(ctx, flight)
		if err != nil {
			return err
		}
	}

	// 更新預訂
	err = s.bookingRepo.UpdateBooking(ctx, booking)
	if err != nil {
		return err
	}

	// 重新評估風險
	riskScore, err := s.overbookingService.AssessRisk(ctx, booking)
	if err != nil {
		logger.Error("Failed to reassess booking risk", zap.Error(err), zap.Int("bookingID", booking.ID))
	} else {
		booking.RiskScore = riskScore
		s.bookingRepo.UpdateBooking(ctx, booking)
	}

	// 發送更新通知
	s.notifyService.NotifyPassenger(ctx, booking, "Your booking has been updated.")

	return nil
}

func (s *bookingService) CancelBooking(ctx context.Context, bookingID int) error {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	flight, err := s.flightRepo.GetFlightByID(ctx, booking.FlightID)
	if err != nil {
		return err
	}

	// 更新航班座位信息
	switch booking.Class {
	case "economy":
		flight.EconomySeats.Booked--
	case "business":
		flight.BusinessSeats.Booked--
	case "first":
		flight.FirstClassSeats.Booked--
	}

	err = s.flightRepo.UpdateFlight(ctx, flight)
	if err != nil {
		return err
	}

	// 取消預訂
	booking.Status = "cancelled"
	err = s.bookingRepo.UpdateBooking(ctx, booking)
	if err != nil {
		return err
	}

	// 發送取消通知
	s.notifyService.NotifyPassenger(ctx, booking, "Your booking has been cancelled.")

	return nil
}

func (s *bookingService) ListBookingsByPassenger(ctx context.Context, passengerID int) ([]*models.Booking, error) {
	return s.bookingRepo.GetBookingsByPassengerID(ctx, passengerID)
}

func (s *bookingService) CheckIn(ctx context.Context, bookingID int) error {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if booking.Status != "confirmed" {
		return errors.New("booking is not in a confirmed state")
	}

	booking.HasCheckedIn = true
	booking.CheckInTime = time.Now()

	err = s.bookingRepo.UpdateBooking(ctx, booking)
	if err != nil {
		return err
	}

	// 重新評估風險
	riskScore, err := s.overbookingService.AssessRisk(ctx, booking)
	if err != nil {
		logger.Error("Failed to reassess booking risk after check-in", zap.Error(err), zap.Int("bookingID", booking.ID))
	} else {
		booking.RiskScore = riskScore
		s.bookingRepo.UpdateBooking(ctx, booking)
	}

	// 發送登機牌
	s.notifyService.NotifyPassenger(ctx, booking, "You have successfully checked in. Here is your boarding pass.")

	return nil
}

func (s *bookingService) calculateAvailableSeats(flight *models.Flight, class string) int {
	switch class {
	case "economy":
		return int(float64(flight.EconomySeats.Total)*(1+flight.EconomySeats.OverbookingRatio)) - flight.EconomySeats.Booked
	case "business":
		return int(float64(flight.BusinessSeats.Total)*(1+flight.BusinessSeats.OverbookingRatio)) - flight.BusinessSeats.Booked
	case "first":
		return int(float64(flight.FirstClassSeats.Total)*(1+flight.FirstClassSeats.OverbookingRatio)) - flight.FirstClassSeats.Booked
	default:
		return 0
	}
}
