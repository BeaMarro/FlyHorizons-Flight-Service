package repositories_test

import (
	"flyhorizons-flightservice/repositories"
	entities "flyhorizons-flightservice/repositories/entity"
	"log"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Create a test version of FlightRepository that uses an in-memory SQLite database
type TestFlightRepository struct {
	repositories.BaseRepository
}

func (repo *TestFlightRepository) CreateConnection() (*gorm.DB, error) {
	if repo.DB != nil {
		return repo.DB, nil
	}

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to create SQLite database: %v", err)
		return nil, err
	}

	// Auto migrate entities for the test database
	err = db.AutoMigrate(&entities.FlightEntity{})
	if err != nil {
		return nil, err
	}

	repo.DB = db
	return db, nil
}

func NewTestFlightRepository() *repositories.FlightRepository {
	baseRepo := &TestFlightRepository{}
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{}) // No shared cache
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Auto-migrate tables for the test database
	if err := db.AutoMigrate(&entities.FlightEntity{}); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	baseRepo.DB = db
	return repositories.NewFlightRepository(&baseRepo.BaseRepository)
}

// Adds flights to the database on every run
func setupFlights(repo *repositories.FlightRepository) []entities.FlightEntity {
	// Flights
	testFlights := []entities.FlightEntity{
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

	// Add flights to the test database
	for _, flight := range testFlights {
		createdFlight := repo.Create(flight)
		log.Printf("Created flight: %+v", createdFlight)
	}

	return testFlights
}

// Integration Database Tests
func TestFlightRepositoryGetAllReturnsFlights(t *testing.T) {
	// Arrange
	flightRepo := NewTestFlightRepository()
	testFlights := setupFlights(flightRepo)

	// Act
	flights := flightRepo.GetAll()

	// Assert
	assert.Equal(t, testFlights, flights)
}

func TestFlightRepositoryGetByValidFlightCodeReturnsFlight(t *testing.T) {
	// Arrange
	flightRepo := NewTestFlightRepository()
	testFlights := setupFlights(flightRepo)
	flightCode := "FR788"

	// Act
	flight := flightRepo.GetByFlightCode(flightCode)

	// Assert
	assert.Equal(t, testFlights[0], flight)
}

func TestFlightRepositoryGetByInvalidFlightCodeReturnsEmptyFlight(t *testing.T) {
	// Arrange
	flightRepo := NewTestFlightRepository()
	invalidFlightCode := "FR999"

	// Act
	flight := flightRepo.GetByFlightCode(invalidFlightCode)

	// Assert
	assert.Equal(t, entities.FlightEntity{}, flight)
}

func TestFlightRepositoryCreateFlightReturnsNewFlight(t *testing.T) {
	// Arrange
	flightRepo := NewTestFlightRepository()
	testFlights := setupFlights(flightRepo)
	flightEntity := entities.FlightEntity{
		FlightCode:        "FR750",
		Departure:         "BLQ",
		Arrival:           "FCO",
		DurationInMinutes: 40,
		DepartureTime:     time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
		DepartureDays:     "[1, 3]",
		CreatedAt:         time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
	}

	// Act
	flight := flightRepo.Create(flightEntity)
	flights := flightRepo.GetAll()

	// Assert
	assert.Len(t, flights, len(testFlights)+1)
	assert.Equal(t, flightEntity, flight)
}

func TestDeleteByValidFlightCodeReturnsTrue(t *testing.T) {
	// Arrange
	flightRepo := NewTestFlightRepository()
	testFlights := setupFlights(flightRepo)
	flightCode := "FR788"

	// Act
	isDeleted := flightRepo.DeleteByFlightCode(flightCode)
	flights := flightRepo.GetAll()

	// Assert
	assert.Len(t, flights, len(testFlights)-1)
	assert.True(t, isDeleted)
}

func TestDeleteByInvalidFlightCodeReturnsFalse(t *testing.T) {
	// Arrange
	flightRepo := NewTestFlightRepository()
	testFlights := setupFlights(flightRepo)
	invalidFlightCode := "FR7999"

	// Act
	isDeleted := flightRepo.DeleteByFlightCode(invalidFlightCode)
	flights := flightRepo.GetAll()

	// Assert
	assert.Len(t, flights, len(testFlights))
	assert.False(t, isDeleted)
}

func TestUpdateValidFlightReturnsUpdatedFlight(t *testing.T) {
	// Arrange
	flightRepo := NewTestFlightRepository()
	testFlights := setupFlights(flightRepo)
	// Update all flight fields
	updatedFlight := entities.FlightEntity{
		FlightCode:        "FR788",
		Departure:         "BLQ",
		Arrival:           "EIN",
		DurationInMinutes: 120,
		DepartureTime:     time.Date(2025, time.March, 30, 15, 30, 0, 0, time.UTC),
		DepartureDays:     "[2, 4]",
		CreatedAt:         time.Date(2025, time.March, 30, 15, 30, 0, 0, time.UTC),
	}

	// Act
	flight := flightRepo.Update(updatedFlight)

	// Assert
	assert.Equal(t, updatedFlight, flight)
	assert.NotNil(t, testFlights)
}
