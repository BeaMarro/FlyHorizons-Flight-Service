package converter_test

import (
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/models/enums"
	entities "flyhorizons-flightservice/repositories/entity"
	"flyhorizons-flightservice/services/converter"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

type TestFlightConverter struct {
}

// Setup
func setup() converter.FlightConverter {
	return converter.FlightConverter{}
}

func getFlightEntity() entities.FlightEntity {
	return entities.FlightEntity{
		FlightCode:        "FR788",
		Departure:         "BLQ",
		Arrival:           "EIN",
		DurationInMinutes: 140,
		DepartureTime:     time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
		DepartureDays:     "[1, 5]",
		CreatedAt:         time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
	}
}

func getFlight() models.Flight {
	return models.Flight{
		FlightCode:        "FR788",
		Departure:         "BLQ",
		Arrival:           "EIN",
		DurationInMinutes: 140,
		DepartureTime:     time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
		DepartureDays:     []enums.Day{enums.Monday, enums.Friday},
	}
}

// Converter Tests
func TestConvertFlightToFlightEntityReturnsFlightEntity(t *testing.T) {
	// Arrange
	flightConverter := setup()

	// Act
	flightEntity := flightConverter.ConvertFlightToFlightEntity(getFlight())

	// Assert
	assert.Equal(t, flightEntity.FlightCode, getFlightEntity().FlightCode)
	assert.Equal(t, flightEntity.Departure, getFlightEntity().Departure)
	assert.Equal(t, flightEntity.Arrival, getFlightEntity().Arrival)
	assert.Equal(t, flightEntity.DurationInMinutes, getFlightEntity().DurationInMinutes)
	assert.Equal(t, flightEntity.DepartureTime, getFlightEntity().DepartureTime)
}

func TestConvertFlightEntityToFlightReturnsFlight(t *testing.T) {
	// Arrange
	flightConverter := setup()

	// Act
	flight := flightConverter.ConvertFlightEntityToFlight(getFlightEntity())

	// Assert
	assert.Equal(t, flight, getFlight())
}
