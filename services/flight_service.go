package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"airline-booking/logger"
	"airline-booking/models"
	"airline-booking/repositories"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type FlightService interface {
	SearchFlights(ctx context.Context, req models.SearchRequest) (string, error)
	GetSearchResults(requestID string) ([]models.Flight, error)
	ProcessSearchRequests()
	GetSearchQueue() <-chan models.SearchRequest
}

type flightService struct {
	repo         repositories.FlightRepository
	redis        *redis.Client
	searchQueue  chan models.SearchRequest
	results      map[string][]models.Flight
	resultsMutex sync.RWMutex
}

func NewFlightService(repo repositories.FlightRepository, redis *redis.Client) FlightService {
	return &flightService{
		repo:        repo,
		redis:       redis,
		searchQueue: make(chan models.SearchRequest, 100), // 緩衝區大小可以根據需求調整
		results:     make(map[string][]models.Flight),
	}
}

func (s *flightService) SearchFlights(ctx context.Context, req models.SearchRequest) (string, error) {
	requestID := fmt.Sprintf("%s-%s-%s-%d", req.Origin, req.Destination, req.Date.Format("2006-01-02"), time.Now().UnixNano())

	s.searchQueue <- req

	logger.LogWithTracing(ctx, "Search request queued",
		zap.String("requestID", requestID),
		zap.String("origin", req.Origin),
		zap.String("destination", req.Destination))

	return requestID, nil
}

func (s *flightService) GetSearchResults(requestID string) ([]models.Flight, error) {
	s.resultsMutex.RLock()
	defer s.resultsMutex.RUnlock()

	if flights, ok := s.results[requestID]; ok {
		return flights, nil
	}

	return nil, fmt.Errorf("results not ready or not found")
}

func (s *flightService) ProcessSearchRequests() {
	for req := range s.searchQueue {
		ctx := context.Background()
		if err := s.processRequest(ctx, req); err != nil {
			logger.Error("Failed to process search request",
				zap.Error(err),
				zap.String("origin", req.Origin),
				zap.String("destination", req.Destination))
		}
	}
}

func (s *flightService) processRequest(ctx context.Context, req models.SearchRequest) error {
	cacheKey := fmt.Sprintf("flights:%s:%s:%s:%d:%d",
		req.Origin, req.Destination, req.Date.Format("2006-01-02"), req.Page, req.PageSize)

	cachedData, err := s.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var flights []models.Flight
		err = json.Unmarshal(cachedData, &flights)
		if err == nil {
			s.storeResults(cacheKey, flights)
			return nil
		}
	}

	// 如果緩存中沒有，則從數據庫中搜索
	flights, err := s.repo.SearchFlights(ctx, req)
	if err != nil {
		return err
	}

	// 將結果存入 Redis 緩存
	cacheData, err := json.Marshal(flights)
	if err == nil {
		// 使用帶有 context 的 Set 方法
		// 好處：
		// 1. 可以傳遞超時設置，允許控制操作的最長執行時間
		// 2. 支持取消操作，如果上層調用被取消，Redis 操作也會被取消
		// 3. 可以傳遞請求級別的元數據，如追蹤 ID，有助於分佈式追蹤
		// 4. 提高了代碼的一致性，與其他使用 context 的 Go 標準庫和第三方庫保持一致
		if err := s.redis.Set(ctx, cacheKey, cacheData, 15*time.Minute).Err(); err != nil {
			logger.Error("Failed to set cache", zap.Error(err), zap.String("cacheKey", cacheKey))
		}
	}

	s.storeResults(cacheKey, flights)
	return nil
}

func (s *flightService) storeResults(requestID string, flights []models.Flight) {
	s.resultsMutex.Lock()
	defer s.resultsMutex.Unlock()
	s.results[requestID] = flights
}

func (s *flightService) GetSearchQueue() <-chan models.SearchRequest {
	return s.searchQueue
}

func (s *flightService) calculateAvailableSeats(flight *models.Flight, class string) int {
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

func (s *flightService) BookFlight(ctx context.Context, flightID int, class string, numSeats int) error {
	flight, err := s.repo.GetFlightByID(ctx, flightID)
	if err != nil {
		return err
	}

	availableSeats := s.calculateAvailableSeats(flight, class)
	if numSeats > availableSeats {
		return errors.New("not enough seats available")
	}

	// 執行預訂邏輯
	// ...

	return nil
}
