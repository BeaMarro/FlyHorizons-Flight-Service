package endtoend

import (
	"bytes"
	"encoding/json"
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/models/enums"
	"flyhorizons-flightservice/repositories"
	entities "flyhorizons-flightservice/repositories/entity"
	"flyhorizons-flightservice/routes"
	"flyhorizons-flightservice/services"
	"flyhorizons-flightservice/services/converter"
	mock_repositories "flyhorizons-flightservice/tests/mocks"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type FlightServiceEndToEndTests struct {
	repositories.BaseRepository
}

// Create a test version of BaseRepository that uses an in-memory SQLite database
func (repo *FlightServiceEndToEndTests) CreateConnection() (*gorm.DB, error) {
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
	baseRepo := &FlightServiceEndToEndTests{}
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
func setupFlights(repo *repositories.FlightRepository) {
	// Users
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

	// Add users to the test database
	for _, flight := range testFlights {
		createdUser := repo.Create(flight)
		log.Printf("Created flight: %+v", createdUser)
	}
}

// Setup
func setupFlightService(repo *repositories.FlightRepository) *services.FlightService {
	flightConverter := converter.FlightConverter{}
	return services.NewFlightService(repo, flightConverter)
}

func setupFlightRouter(service services.FlightService, gatewayAuthMiddleware *mock_repositories.MockGatewayAuthMiddleware) *gin.Engine {
	router := gin.Default()
	routes.RegisterFlightRoutes(router, &service, gatewayAuthMiddleware)
	return router
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

// End-to-End Tests
func TestEndToEndGetAllReturnsFlights(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	flightService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	mockFlights := getFlights()
	// Setup router
	router := setupFlightRouter(*flightService, mockAPIGatewayMiddleware)

	url := "/flights"
	httpRequest, _ := http.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// Unmarshal the JSON response
	var flights []models.Flight
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &flights)

	assert.NoError(t, err)
	assert.Equal(t, mockFlights, flights)
}

func TestEndToEndGetFlightByExistingFlightCodeReturnsFlight(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	flightService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	mockFlight := getFlights()[0]
	mockFlightCode := mockFlight.FlightCode
	// Setup router
	router := setupFlightRouter(*flightService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/flights/%s", mockFlightCode)
	httpRequest, _ := http.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// Unmarshal the JSON response
	var flight models.Flight
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &flight)

	assert.NoError(t, err)
	assert.Equal(t, mockFlight, flight)
}

func TestEndToEndGetFlightByNonExistingFlightCodeReturnsNotFoundError(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	flightService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	mockFlightCode := "FR799"
	// Setup router
	router := setupFlightRouter(*flightService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/flights/%s", mockFlightCode)
	httpRequest, _ := http.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestEndToEndCreateNonExistingFlightAsAdminReturnsCreatedFlight(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	flightService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockFlight := models.Flight{
		FlightCode:        "FR787",
		Departure:         "BLQ",
		Arrival:           "EIN",
		DurationInMinutes: 115,
		DepartureTime:     time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
		DepartureDays:     []enums.Day{enums.Monday, enums.Friday},
	}
	// Setup router
	router := setupFlightRouter(*flightService, mockAPIGatewayMiddleware)
	// Make the JSON to create the flight
	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("POST", "/flights/", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                           // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	var flight models.Flight
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &flight)
	assert.NoError(t, err)
	assert.Equal(t, mockFlight.FlightCode, flight.FlightCode)
	assert.Equal(t, mockFlight.Departure, flight.Departure)
	assert.Equal(t, mockFlight.Arrival, flight.Arrival)
	assert.Equal(t, mockFlight.DurationInMinutes, flight.DurationInMinutes)
	assert.Equal(t, mockFlight.DepartureTime, flight.DepartureTime)
	assert.Equal(t, mockFlight.DepartureDays, flight.DepartureDays)
}

func TestEndToEndCreateExistingFlightAsAdminReturnsConflictError(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	flightService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockFlight := getFlights()[0]
	// Setup router
	router := setupFlightRouter(*flightService, mockAPIGatewayMiddleware)
	// Make the JSON to create the flight
	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("POST", "/flights/", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                           // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusConflict, responseRecorder.Code)
}

func TestEndToEndCreateFlightAsUserReturnsAccessDenied(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	flightService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	mockFlight := getFlights()[0]
	// Setup router
	router := setupFlightRouter(*flightService, mockAPIGatewayMiddleware)
	// Make the JSON to create the flight
	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("POST", "/flights/", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                           // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestEndToEndDeleteExistingFlightAsAdminReturnsDeletedSuccessfully(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	userService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	bearerToken := "Bearer mocktoken12345"
	mockFlightCode := getFlights()[0].FlightCode
	// Setup router
	router := setupFlightRouter(*userService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/flights/%s", mockFlightCode)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestDeleteNonExistingFlightAsAdminReturnsNotFoundError(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	userService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	bearerToken := "Bearer mocktoken12345"
	mockFlightCode := "FR787"
	// Setup router
	router := setupFlightRouter(*userService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/flights/%s", mockFlightCode)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestDeleteFlightAsUserReturnsAccessDenied(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	userService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	bearerToken := "Bearer mocktoken12345"
	mockFlightCode := "FR787"
	// Setup router
	router := setupFlightRouter(*userService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/flights/%s", mockFlightCode)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestUpdateExistingFlightAsAdminReturnsUpdateFlight(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	flightService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockFlight := models.Flight{
		FlightCode:        "FR788",
		Departure:         "BLQ",
		Arrival:           "EIN",
		DurationInMinutes: 115,
		DepartureTime:     time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
		DepartureDays:     []enums.Day{enums.Monday, enums.Friday},
	}
	// Setup router
	router := setupFlightRouter(*flightService, mockAPIGatewayMiddleware)
	// Make the JSON to create the flight
	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("PUT", "/flights/", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                          // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var flight models.Flight
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &flight)
	assert.NoError(t, err)
	assert.Equal(t, mockFlight.FlightCode, flight.FlightCode)
	assert.Equal(t, mockFlight.Departure, flight.Departure)
	assert.Equal(t, mockFlight.Arrival, flight.Arrival)
	assert.Equal(t, mockFlight.DurationInMinutes, flight.DurationInMinutes)
	assert.Equal(t, mockFlight.DepartureTime, flight.DepartureTime)
	assert.Equal(t, mockFlight.DepartureDays, flight.DepartureDays)
}

func TestUpdateNonExistingFlightAsAdminReturnsNotFoundError(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	flightService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockFlight := models.Flight{
		FlightCode:        "FR799",
		Departure:         "BLQ",
		Arrival:           "EIN",
		DurationInMinutes: 115,
		DepartureTime:     time.Date(2025, time.April, 1, 15, 30, 0, 0, time.UTC),
		DepartureDays:     []enums.Day{enums.Monday, enums.Friday},
	}
	// Setup router
	router := setupFlightRouter(*flightService, mockAPIGatewayMiddleware)
	// Make the JSON to create the flight
	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("PUT", "/flights/", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                          // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestUpdateFlightAsUserReturnsAccessDenied(t *testing.T) {
	// Arrange
	// Setup repository
	flightRepo := NewTestFlightRepository()
	setupFlights(flightRepo)
	// Setup service
	flightService := setupFlightService(flightRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	mockFlight := getFlights()[0]
	// Setup router
	router := setupFlightRouter(*flightService, mockAPIGatewayMiddleware)
	// Make the JSON to create the flight
	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("PUT", "/flights/", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                          // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}
