package services_test

import (
	"context"
	"testing"
	"time"

	"airline-booking/logger"
	"airline-booking/mocks"
	"airline-booking/models"
	"airline-booking/services"

	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	testLogger, _ := zap.NewDevelopment()
	logger.SetLoggerForTest(testLogger)
}

func TestFlightService_SearchFlights(t *testing.T) {
	// 為每個測試設置一個新的 logger
	testLogger, _ := zap.NewDevelopment()
	logger.SetLoggerForTest(testLogger)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFlightRepository(ctrl)
	mockRedis, _ := redismock.NewClientMock()

	// 設置預期行為
	mockRepo.EXPECT().SearchFlights(gomock.Any(), gomock.Any()).Return([]models.Flight{}, nil).AnyTimes()

	service := services.NewFlightService(mockRepo, mockRedis)

	ctx := context.Background()
	req := models.SearchRequest{
		Origin:      "New York",
		Destination: "London",
		Date:        time.Now(),
	}

	requestID, err := service.SearchFlights(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, requestID)

	// 驗證 requestID 格式
	assert.Regexp(t, `^New York-London-\d{4}-\d{2}-\d{2}-\d+$`, requestID)

	// 驗證請求被加入隊列
	select {
	case receivedReq := <-service.GetSearchQueue():
		assert.Equal(t, req.Origin, receivedReq.Origin)
		assert.Equal(t, req.Destination, receivedReq.Destination)
	case <-time.After(time.Second):
		t.Error("Request was not added to the queue within the expected time")
	}
}
