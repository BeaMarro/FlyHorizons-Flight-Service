package services_test

import (
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/models/enums"
	"flyhorizons-flightservice/services"
	strategies "flyhorizons-flightservice/services/sort_strategies"
	"log"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

type FilterFlightServiceTest struct {
}

// Setup
func setupFlightFilterService() services.FlightFilterService {
	flightFilterService := services.FlightFilterService{}
	return flightFilterService
}

// Service Unit Tests
func TestFilterByAllReturnsFilteredFlights(t *testing.T) {
	// Arrange
	// Filter Data
	departureAirport := "EIN"
	arrivalAirport := "BLQ"
	departureDate := time.Date(2025, time.March, 14, 0, 0, 0, 0, time.Local)
	returnDate := time.Date(2025, time.March, 14, 0, 0, 0, 0, time.Local)

	// Convert to pointer (nullable)
	departureAirportPtr := &departureAirport
	arrivalAirportPtr := &arrivalAirport
	departureDatePtr := &departureDate
	returnDatePtr := &returnDate

	flightFilterService := setupFlightFilterService()

	flights := []models.Flight{
		{
			FlightCode:        "FR788",
			Departure:         "BLQ",
			Arrival:           "EIN",
			DurationInMinutes: 140,
			DepartureTime:     time.Now(),
			DepartureDays:     []enums.Day{enums.Friday},
		},
		{
			FlightCode:        "FR789",
			Departure:         "EIN",
			Arrival:           "BLQ",
			DurationInMinutes: 120,
			DepartureTime:     time.Now(),
			DepartureDays:     []enums.Day{enums.Friday},
		},
	}

	expected := []models.Flight{
		{
			FlightCode:        "FR789",
			Departure:         "EIN",
			Arrival:           "BLQ",
			DurationInMinutes: 120,
			DepartureTime:     time.Now(),
			DepartureDays:     []enums.Day{enums.Friday},
		},
	}
	// Setup strategies
	arrivalStrategy := strategies.ArrivalAirportStrategy{}
	departureStrategy := strategies.DepartureAirportStrategy{}
	dateStrategy := strategies.DateRangeStrategy{}

	flightFilterService.AddStrategy(arrivalStrategy)
	flightFilterService.AddStrategy(departureStrategy)
	flightFilterService.AddStrategy(dateStrategy)

	// Act
	start := time.Now()
	filteredFlights := flightFilterService.Filter(flights, departureAirportPtr, arrivalAirportPtr, departureDatePtr, returnDatePtr)
	elapsed := time.Since(start)

	// Assert
	assert.Equal(t, expected, filteredFlights)
	log.Printf("Execution time: %s\n", elapsed)
}

func TestFilterByDepartureAndArrivalReturnsFlights(t *testing.T) {
	// Arrange
	// Filter Data
	departureAirport := "BLQ"
	arrivalAirport := "EIN"

	// Convert to pointer (nullable)
	departureAirportPtr := &departureAirport
	arrivalAirportPtr := &arrivalAirport

	flightFilterService := setupFlightFilterService()

	flights := []models.Flight{
		{
			FlightCode:        "FR788",
			Departure:         "BLQ",
			Arrival:           "EIN",
			DurationInMinutes: 140,
			DepartureTime:     time.Now(),
			DepartureDays:     []enums.Day{enums.Friday},
		},
		{
			FlightCode:        "FR789",
			Departure:         "EIN",
			Arrival:           "BLQ",
			DurationInMinutes: 120,
			DepartureTime:     time.Now(),
			DepartureDays:     []enums.Day{enums.Friday},
		},
	}

	expected := []models.Flight{
		{
			FlightCode:        "FR788",
			Departure:         "BLQ",
			Arrival:           "EIN",
			DurationInMinutes: 140,
			DepartureTime:     time.Now(),
			DepartureDays:     []enums.Day{enums.Friday},
		},
	}
	// Setup strategies
	arrivalStrategy := strategies.ArrivalAirportStrategy{}
	departureStrategy := strategies.DepartureAirportStrategy{}

	flightFilterService.AddStrategy(arrivalStrategy)
	flightFilterService.AddStrategy(departureStrategy)

	// Act
	filteredFlights := flightFilterService.Filter(flights, departureAirportPtr, arrivalAirportPtr, nil, nil)

	// Assert
	assert.Equal(t, expected, filteredFlights)
}
