package routes

import (
	"flyhorizons-flightservice/services"
	"flyhorizons-flightservice/services/interfaces"
	strategies "flyhorizons-flightservice/services/sort_strategies"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handles the flight filtering functionality
func RegisterFilterFlightRoutes(router *gin.Engine, flightService interfaces.FlightService) {
	router.GET("/flights/filter", func(ctx *gin.Context) {
		var departureAirport *string
		var arrivalAirport *string

		// Extracting the query parameters (if present)
		if departure := ctx.DefaultQuery("departureAirport", ""); departure != "" {
			departureAirport = &departure
		}

		if arrival := ctx.DefaultQuery("arrivalAirport", ""); arrival != "" {
			arrivalAirport = &arrival
		}

		departureDateStr := ctx.DefaultQuery("departureDate", "")
		returnDateStr := ctx.DefaultQuery("returnDate", "")

		var departureDate, returnDate *time.Time

		// Parse departure and return dates as datetime objects
		if departureDateStr != "" {
			parsedDepartureDate, err := time.Parse("2006-01-02", departureDateStr)
			if err == nil {
				departureDate = &parsedDepartureDate
			}
		}

		if returnDateStr != "" {
			parsedReturnDate, err := time.Parse("2006-01-02", returnDateStr)
			if err == nil {
				returnDate = &parsedReturnDate
			}
		}

		flightFilterService := services.FlightFilterService{}

		// Add strategies based on query parameters
		if arrivalAirport != nil {
			arrivalStrategy := strategies.ArrivalAirportStrategy{}
			flightFilterService.AddStrategy(arrivalStrategy)
		}
		if departureAirport != nil {
			departureStrategy := strategies.DepartureAirportStrategy{}
			flightFilterService.AddStrategy(departureStrategy)
		}
		if departureDate != nil || returnDate != nil {
			dateStrategy := strategies.DateRangeStrategy{}
			flightFilterService.AddStrategy(dateStrategy)
		}

		// Get all flights from the flightService
		flights := flightService.GetAll(ctx.Request.Context())

		// Filter the flights using the filter service and the query parameters (if applicable)
		filteredFlights := flightFilterService.Filter(flights, departureAirport, arrivalAirport, departureDate, returnDate)

		if len(filteredFlights) > 0 {
			ctx.JSON(http.StatusOK, filteredFlights)
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "No flights found matching the criteria"})
		}
	})
}
