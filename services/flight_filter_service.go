package services

import (
	"flyhorizons-flightservice/models"
	"time"
)

type FlightFilterService struct {
	Strategies []FilterStrategy
}

func (service *FlightFilterService) AddStrategy(strategy FilterStrategy) {
	service.Strategies = append(service.Strategies, strategy)
}

func (service *FlightFilterService) Filter(flights []models.Flight, departureAirport *string, arrivalAirport *string, departureDate *time.Time, returnDate *time.Time) []models.Flight {
	// Applies strategies in a sequence
	for _, strategy := range service.Strategies {
		flights = strategy.Filter(flights, departureAirport, arrivalAirport, departureDate, returnDate)
	}
	return flights
}
