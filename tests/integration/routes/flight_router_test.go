package routes_test

import (
	"bytes"
	"encoding/json"
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/models/enums"
	"flyhorizons-flightservice/routes"
	"flyhorizons-flightservice/services/errors"
	mock_repositories "flyhorizons-flightservice/tests/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestFlightRoute struct {
}

// Setup
func setupFlightRouter(mockService *mock_repositories.MockFlightService, gatewayAuthMiddleware *mock_repositories.MockGatewayAuthMiddleware) *gin.Engine {
	router := gin.Default()

	routes.RegisterFlightRoutes(router, mockService, gatewayAuthMiddleware)

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

// Router Integration Tests
func TestGetAllReturnsFlightsJSON(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := new(mock_repositories.MockGatewayAuthMiddleware)
	mockFlights := getFlights()

	mockService.On("GetAll").Return(mockFlights, nil)

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	httpRequest, _ := http.NewRequest("GET", "/flights", nil)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var flights []models.Flight
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &flights)
	assert.NoError(t, err)
	assert.Equal(t, mockFlights, flights)
	mockService.AssertExpectations(t)
}

func TestGetByExistingFlightReturnsFlightJSON(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := new(mock_repositories.MockGatewayAuthMiddleware)
	mockFlight := getFlights()[0]
	mockFlightCode := mockFlight.FlightCode

	mockService.On("GetByFlightCode", mockFlightCode).Return(&mockFlight, nil)

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/flights/%s", mockFlightCode)
	httpRequest, _ := http.NewRequest("GET", url, nil)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var flight models.Flight
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &flight)
	assert.NoError(t, err)
	assert.Equal(t, mockFlight, flight)
	mockService.AssertExpectations(t)
}

func TestGetByNonExistingFlightReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := new(mock_repositories.MockGatewayAuthMiddleware)
	mockFlightCode := "FH9999"

	mockService.On("GetByFlightCode", mockFlightCode).Return(nil, errors.NewFlightNotFoundError(mockFlightCode, 404))

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/flights/%s", mockFlightCode)
	httpRequest, _ := http.NewRequest("GET", url, nil)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)

	var flight models.Flight
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &flight)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestCreateNonExistingFlightAsAdminReturnsCreatedFlight(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockFlight := getFlights()[0]
	mockService.On("Create", mockFlight).Return(&mockFlight, nil)
	bearerToken := "Bearer mocktoken12345"

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the flight
	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("POST", "/flights/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	var flight models.Flight
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &flight)
	assert.NoError(t, err)
	assert.Equal(t, mockFlight, flight)
	mockService.AssertExpectations(t)
}

func TestCreateExistingFlightAsAdminReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockFlight := getFlights()[0]
	mockService.On("Create", mockFlight).Return(nil, errors.NewFlightExistsError(mockFlight.FlightCode, 409))
	bearerToken := "Bearer mocktoken12345"

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("POST", "/flights/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusConflict, responseRecorder.Code)

	var errResponse map[string]interface{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestCreateFlightAsNonAdminRoleReturnsAccessDenied(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	mockFlight := getFlights()[0]
	bearerToken := "Bearer mocktoken12345"

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the airport
	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("POST", "/flights/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)

	var errResponse map[string]interface{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestDeleteExistingFlightAsAdminReturnsHTTPStatusOK(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockFlight := getFlights()[0]
	mockFlightCode := mockFlight.FlightCode
	bearerToken := "Bearer mocktoken12345"
	mockService.On("DeleteByFlightCode", mockFlightCode).Return(true, nil)

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/flights/%s", mockFlightCode)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteNonExistingFlightAsAdminReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockFlightCode := "FH9999"
	bearerToken := "Bearer mocktoken12345"
	mockService.On("DeleteByFlightCode", mockFlightCode).Return(false, errors.NewFlightNotFoundError(mockFlightCode, 404))

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/flights/%s", mockFlightCode)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteFlightAsNonAdminRoleReturnsAccessDenied(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	mockFlightCode := "FH9999"
	bearerToken := "Bearer mocktoken12345"

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the airport
	url := fmt.Sprintf("/flights/%s", mockFlightCode)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)

	var errResponse map[string]interface{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestUpdateExistingFlightAsAdminReturnsUpdatedFlight(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockFlight := getFlights()[0]
	mockService.On("Update", mockFlight).Return(&mockFlight, nil)
	bearerToken := "Bearer mocktoken12345"

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("PUT", "/flights/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var flight models.Flight
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &flight)
	assert.NoError(t, err)
	assert.Equal(t, mockFlight, flight)
	mockService.AssertExpectations(t)
}

func TestUpdateNonExistingFlightAsAdminReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockFlight := getFlights()[0]
	mockService.On("Update", mockFlight).Return(nil, errors.NewFlightNotFoundError(mockFlight.FlightCode, 404))
	bearerToken := "Bearer mocktoken12345"

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("PUT", "/flights/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)

	var errResponse map[string]interface{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestUpdateFlightAsNonAdminRoleReturnsAccessDenied(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	mockFlight := getFlights()[0]
	bearerToken := "Bearer mocktoken12345"

	router := setupFlightRouter(mockService, mockAPIGatewayMiddleware)

	requestBody, _ := json.Marshal(mockFlight)
	httpRequest, _ := http.NewRequest("PUT", "/flights/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)

	var errResponse map[string]interface{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}
