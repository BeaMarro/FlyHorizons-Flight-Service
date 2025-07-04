package strategies

import (
	"flyhorizons-flightservice/models"
	"time"
)

type DepartureAirportStrategy struct{}

func (strategy DepartureAirportStrategy) Filter(flights []models.Flight, depatureAirport *string, arrivalAirport *string, departureDate *time.Time, returnDate *time.Time) []models.Flight {
	var filteredFlights []models.Flight

	if depatureAirport != nil {
		flightDeparture := *depatureAirport
		for _, flight := range flights {
			if flight.Departure == flightDeparture {
				filteredFlights = append(filteredFlights, flight)
			}
		}
	}
	return filteredFlights
}
