package services

import (
	"context"
	"encoding/json"
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/services/converter"
	"flyhorizons-flightservice/services/errors"
	"flyhorizons-flightservice/services/interfaces"
	"time"

	"github.com/redis/go-redis/v9"
)

type FlightService struct {
	flightRepo      interfaces.FlightRepository
	flightConverter converter.FlightConverter
	redisClient     *redis.Client
}

func NewFlightService(repo interfaces.FlightRepository, flightConverter converter.FlightConverter, redisClient *redis.Client) *FlightService {
	return &FlightService{
		flightRepo:      repo,
		flightConverter: flightConverter,
		redisClient:     redisClient,
	}
}

func (flightService *FlightService) GetAll(ctx context.Context) []models.Flight {
	cacheKey := "flights:all"
	cached, err := flightService.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var flights []models.Flight
		if err := json.Unmarshal([]byte(cached), &flights); err == nil {
			return flights
		}
	}

	flightEntities := flightService.flightRepo.GetAll()
	var flights []models.Flight
	for _, flightEntity := range flightEntities {
		flight := flightService.flightConverter.ConvertFlightEntityToFlight(flightEntity)
		flights = append(flights, flight)
	}

	data, err := json.Marshal(flights)
	if err == nil {
		flightService.redisClient.Set(ctx, cacheKey, data, 2*time.Minute)
	}

	return flights
}

func (flightService *FlightService) GetByFlightCode(ctx context.Context, flightCode string) (*models.Flight, error) {
	cacheKey := "flight:" + flightCode
	cached, err := flightService.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var flight models.Flight
		if err := json.Unmarshal([]byte(cached), &flight); err == nil {
			return &flight, nil
		}
	}
	flightEntity := flightService.flightRepo.GetByFlightCode(flightCode)
	flight := flightService.flightConverter.ConvertFlightEntityToFlight(flightEntity)
	data, err := json.Marshal(flightEntity)
	if err == nil {
		flightService.redisClient.Set(ctx, cacheKey, data, 5*time.Minute)
	}
	return &flight, nil
}

func (FlightService *FlightService) FlightExists(ctx context.Context, flightCode string) bool {
	for _, flight := range FlightService.GetAll(ctx) {
		if flight.FlightCode == flightCode {
			return true
		}
	}
	return false
}

func (flightService *FlightService) Create(ctx context.Context, flight models.Flight) (*models.Flight, error) {
	if flightService.FlightExists(ctx, flight.FlightCode) {
		return nil, errors.NewFlightExistsError(flight.FlightCode, 409)
	}
	flightEntity := flightService.flightConverter.ConvertFlightToFlightEntity(flight)
	createdFlightEntity := flightService.flightRepo.Create(flightEntity)
	createdFlight := flightService.flightConverter.ConvertFlightEntityToFlight(createdFlightEntity)

	// Invalidate both single flight and list cache
	flightService.redisClient.Del(ctx, "flight:"+flight.FlightCode)
	flightService.redisClient.Del(ctx, "flights:all")

	return &createdFlight, nil
}

func (flightService *FlightService) DeleteByFlightCode(ctx context.Context, flightCode string) (bool, error) {
	if !flightService.FlightExists(ctx, flightCode) {
		return false, errors.NewFlightNotFoundError(flightCode, 404)
	}
	success := flightService.flightRepo.DeleteByFlightCode(flightCode)

	// Invalidate both single flight and list cache
	flightService.redisClient.Del(ctx, "flight:"+flightCode)
	flightService.redisClient.Del(ctx, "flights:all")

	return success, nil
}

func (flightService *FlightService) Update(ctx context.Context, flight models.Flight) (*models.Flight, error) {
	if !flightService.FlightExists(ctx, flight.FlightCode) {
		return nil, errors.NewFlightNotFoundError(flight.FlightCode, 404)
	}
	flightEntity := flightService.flightConverter.ConvertFlightToFlightEntity(flight)
	updatedFlightEntity := flightService.flightRepo.Update(flightEntity)
	updatedFlight := flightService.flightConverter.ConvertFlightEntityToFlight(updatedFlightEntity)

	// Invalidate both single flight and list cache
	flightService.redisClient.Del(ctx, "flight:"+flight.FlightCode)
	flightService.redisClient.Del(ctx, "flights:all")

	return &updatedFlight, nil
}
