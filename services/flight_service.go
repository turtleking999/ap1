package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"airline-booking/models"
	"airline-booking/repositories"

	"github.com/go-redis/redis/v8"
)

type FlightService interface {
	SearchFlights(req models.SearchRequest) ([]models.Flight, error)
}

type flightService struct {
	repo  repositories.FlightRepository
	redis *redis.Client
}

func NewFlightService(repo repositories.FlightRepository, redis *redis.Client) FlightService {
	return &flightService{repo: repo, redis: redis}
}

func (s *flightService) SearchFlights(req models.SearchRequest) ([]models.Flight, error) {

	cacheKey := fmt.Sprintf("flights:%s:%s:%s:%d:%d",
		req.Origin, req.Destination, req.Date.Format("2006-01-02"), req.Page, req.PageSize)

	ctx := context.Background()
	cachedData, err := s.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var flights []models.Flight
		err = json.Unmarshal(cachedData, &flights)
		if err == nil {
			return flights, nil
		}
	}

	flights, err := s.repo.SearchFlights(req)
	if err != nil {
		return nil, err
	}

	cacheData, err := json.Marshal(flights)
	if err == nil {
		s.redis.Set(ctx, cacheKey, cacheData, 15*time.Minute)
	}

	return flights, nil
}
