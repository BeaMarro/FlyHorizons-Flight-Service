package services

import (
	"flyhorizons-flightservice/models"
	"time"
)

type FilterStrategy interface {
	Filter(flights []models.Flight, depatureAirport *string, arrivalAirport *string, departureDate *time.Time, returnDate *time.Time) []models.Flight
}
