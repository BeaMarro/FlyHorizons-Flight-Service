package routes_test

import (
	"encoding/json"
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/routes"
	mock_repositories "flyhorizons-flightservice/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestFlightFlightRoute struct {
}

// Setup
func setupFlightFilterRouter(mockService *mock_repositories.MockFlightService) *gin.Engine {
	router := gin.Default()

	routes.RegisterFilterFlightRoutes(router, mockService)

	return router
}

// Router Integration Tests
func TestFilterByAllCriteriaReturnsFilteredFlights(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockFlightService)
	allFlights := getFlights()
	expectedFilteredFlights := []models.Flight{
		allFlights[1],
	}
	mockService.On("GetAll").Return(allFlights, nil)

	router := setupFlightFilterRouter(mockService)

	httpRequest, _ := http.NewRequest("GET", "/flights/filter?departureAirport=EIN&arrivalAirport=BLQ&departureDate=2025-4-1&returnDate=2025-5-2", nil)
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var filteredFlights []models.Flight
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &filteredFlights)
	assert.NoError(t, err)
	assert.Equal(t, expectedFilteredFlights, filteredFlights)
	mockService.AssertExpectations(t)
}
