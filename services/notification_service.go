package services

import (
	"airline-booking/models"
	"context"
)

// NotificationService 定義了發送通知的接口
type NotificationService interface {
	NotifyPassenger(ctx context.Context, booking *models.Booking, message string) error
	SendBookingConfirmation(ctx context.Context, booking *models.Booking) error
	SendCheckInReminder(ctx context.Context, booking *models.Booking) error
	SendBoardingPass(ctx context.Context, booking *models.Booking) error
	SendFlightStatusUpdate(ctx context.Context, booking *models.Booking, status string) error
	SendPromotionalOffer(ctx context.Context, passenger *models.Passenger, offer string) error
}

// notificationService 是 NotificationService 接口的空實現
type notificationService struct {
	// 這裡可以添加任何需要的依賴，比如郵件客戶端、SMS 服務等
}

// NewNotificationService 創建一個新的 NotificationService 實例
func NewNotificationService() NotificationService {
	return &notificationService{}
}

func (s *notificationService) NotifyPassenger(ctx context.Context, booking *models.Booking, message string) error {
	// TODO: 實現發送通用通知的邏輯
	return nil
}

func (s *notificationService) SendBookingConfirmation(ctx context.Context, booking *models.Booking) error {
	// TODO: 實現發送預訂確認的邏輯
	return nil
}

func (s *notificationService) SendCheckInReminder(ctx context.Context, booking *models.Booking) error {
	// TODO: 實現發送登機提醒的邏輯
	return nil
}

func (s *notificationService) SendBoardingPass(ctx context.Context, booking *models.Booking) error {
	// TODO: 實現發送登機牌的邏輯
	return nil
}

func (s *notificationService) SendFlightStatusUpdate(ctx context.Context, booking *models.Booking, status string) error {
	// TODO: 實現發送航班狀態更新的邏輯
	return nil
}

func (s *notificationService) SendPromotionalOffer(ctx context.Context, passenger *models.Passenger, offer string) error {
	// TODO: 實現發送促銷優惠的邏輯
	return nil
}
