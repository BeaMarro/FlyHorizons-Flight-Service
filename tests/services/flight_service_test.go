package services_test

import (
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/models/enums"
	entities "flyhorizons-flightservice/repositories/entity"
	"flyhorizons-flightservice/services"
	"flyhorizons-flightservice/services/converter"
	"flyhorizons-flightservice/services/errors"
	mock_repositories "flyhorizons-flightservice/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestFlightService struct {
}

// Setup
func setupFlightService() (*mock_repositories.MockFlightRepository, *services.FlightService) {
	mockRepo := new(mock_repositories.MockFlightRepository)
	flightConverter := new(converter.FlightConverter)
	flightService := services.NewFlightService(mockRepo, *flightConverter)
	return mockRepo, flightService
}

func getFlightEntities() []entities.FlightEntity {
	return []entities.FlightEntity{
		{
			FlightCode:        "FR788",
			Departure:         "BLQ",
			Arrival:           "EIN",
			DurationInMinutes: 140,
			DepartureTime:     time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
			DepartureDays:     "[1, 5]",
			CreatedAt:         time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
		},
		{
			FlightCode:        "FR789",
			Departure:         "EIN",
			Arrival:           "BLQ",
			DurationInMinutes: 120,
			DepartureTime:     time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
			DepartureDays:     "[1, 3]",
			CreatedAt:         time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
		},
	}
}

func getFlights() []models.Flight {
	return []models.Flight{
		{
			FlightCode:        "FR788",
			Departure:         "BLQ",
			Arrival:           "EIN",
			DurationInMinutes: 140,
			DepartureTime:     time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
			DepartureDays:     []enums.Day{enums.Monday, enums.Friday},
		},
		{
			FlightCode:        "FR789",
			Departure:         "EIN",
			Arrival:           "BLQ",
			DurationInMinutes: 120,
			DepartureTime:     time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
			DepartureDays:     []enums.Day{enums.Monday, enums.Wednesday},
		},
	}
}

// Service Unit Tests
func TestGetAllReturnsFlights(t *testing.T) {
	// Arrange
	mockRepo, flightService := setupFlightService()
	mockRepo.On("GetAll").Return(getFlightEntities())

	// Act
	all_flights := flightService.GetAll()

	// Assert
	assert.Equal(t, getFlights(), all_flights)
}

func TestGetByValidFlightCodeReturnsMatchingFlight(t *testing.T) {
	// Arrange
	mockRepo, flightService := setupFlightService()
	mockRepo.On("GetAll").Return(getFlightEntities())
	flightCode := "FR788"
	expectedFlight := getFlights()[0]
	mockRepo.On("GetByFlightCode", flightCode).Return(getFlightEntities()[0], nil)

	// Act
	flight, err := flightService.GetByFlightCode(flightCode)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedFlight, *flight)
}

func TestGetByInvalidFlightCodeThrowsException(t *testing.T) {
	// Arrange
	mockRepo, flightService := setupFlightService()
	flightCode := "FR788"
	errorNotFound := errors.NewFlightNotFoundError(flightCode, 404)
	mockRepo.On("GetByFlightCode", flightCode).Return(entities.FlightEntity{}, errorNotFound)

	// Act
	flight, err := flightService.GetByFlightCode(flightCode)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errorNotFound, err)
	assert.Nil(t, flight)
}

func TestCreateNonExistingFlightReturnsCreatedFlight(t *testing.T) {
	// Arrange
	mockRepo, flightService := setupFlightService()
	mockRepo.On("GetAll").Return([]entities.FlightEntity{})
	flightEntity := getFlightEntities()[0]
	flight := getFlights()[0]
	mockRepo.On("Create", mock.MatchedBy(func(u entities.FlightEntity) bool {
		return u.FlightCode == flightEntity.FlightCode
	})).Return(flightEntity)

	// Act
	createdFlight, err := flightService.Create(flight)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, flight.FlightCode, createdFlight.FlightCode)
}

func TestCreateExistingFlightThrowsException(t *testing.T) {
	// Arrange
	mockRepo, flightService := setupFlightService()
	flight := getFlights()[0]
	flightEntity := getFlightEntities()[0]

	mockRepo.On("GetAll").Return([]entities.FlightEntity{getFlightEntities()[0]})
	mockRepo.On("Create", mock.MatchedBy(func(u entities.FlightEntity) bool {
		return u.FlightCode == flightEntity.FlightCode
	})).Return(flightEntity)

	// Act
	createdFlight, err := flightService.Create(flight)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdFlight)
	assert.Equal(t, errors.NewFlightExistsError(flight.FlightCode, 409), err)
}

func TestDeleteByExistingFlightCodeReturnsTrue(t *testing.T) {
	// Arrange
	mockRepo, flightService := setupFlightService()
	flightCode := getFlightEntities()[0].FlightCode
	mockRepo.On("GetAll").Return([]entities.FlightEntity{getFlightEntities()[0]})
	mockRepo.On("DeleteByFlightCode", flightCode).Return(true)

	// Act
	isDeleted, err := flightService.DeleteByFlightCode(flightCode)

	// Assert
	assert.NoError(t, err)
	assert.True(t, isDeleted)
}

func TestDeleteByNonExistingFlightCodeThrowsException(t *testing.T) {
	// Arrange
	mockRepo, flightService := setupFlightService()
	invalidFlightCode := "FR9999"
	mockRepo.On("GetAll").Return([]entities.FlightEntity{getFlightEntities()[1]})
	mockRepo.On("DeleteByFlightCode", invalidFlightCode).Return(false)

	// Act
	isDeleted, err := flightService.DeleteByFlightCode(invalidFlightCode)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.NewFlightNotFoundError(invalidFlightCode, 404), err)
	assert.False(t, isDeleted)
}

func TestUpdateByExistingFlightReturnsUpdatedFlight(t *testing.T) {
	// Arrange
	mockRepo, flightService := setupFlightService()
	flight := getFlights()[0]
	flightEntity := getFlightEntities()[0]

	mockRepo.On("GetAll").Return(getFlightEntities())
	mockRepo.On("Update", mock.MatchedBy(func(u entities.FlightEntity) bool {
		return u.FlightCode == flightEntity.FlightCode
	})).Return(flightEntity)

	// Act
	updateFlight, err := flightService.Update(flight)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, flight.FlightCode, updateFlight.FlightCode)
	assert.Equal(t, flight.Departure, updateFlight.Departure)
	assert.Equal(t, flight.Arrival, updateFlight.Arrival)
	assert.Equal(t, flight.DurationInMinutes, updateFlight.DurationInMinutes)
	assert.Equal(t, flight.DepartureDays, updateFlight.DepartureDays)
	assert.Equal(t, flight.DepartureTime, updateFlight.DepartureTime)
	assert.Equal(t, flight.BasePrice, updateFlight.BasePrice)
}

func TestUpdateByNonExistingFlightThrowsException(t *testing.T) {
	// Arrange
	mockRepo, flightService := setupFlightService()
	flight := getFlights()[0]
	flightEntity := getFlightEntities()[0]

	mockRepo.On("GetAll").Return([]entities.FlightEntity{})
	mockRepo.On("Update", mock.MatchedBy(func(u entities.FlightEntity) bool {
		return u.FlightCode == flightEntity.FlightCode
	})).Return(flightEntity)

	// Act
	updateFlight, err := flightService.Update(flight)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.NewFlightNotFoundError(flight.FlightCode, 404), err)
	assert.Nil(t, updateFlight)
}
