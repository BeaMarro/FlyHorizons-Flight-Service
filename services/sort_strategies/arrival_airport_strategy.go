package strategies

import (
	"flyhorizons-flightservice/models"
	"time"
)

type ArrivalAirportStrategy struct{}

func (strategy ArrivalAirportStrategy) Filter(flights []models.Flight, depatureAirport *string, arrivalAirport *string, departureDate *time.Time, returnDate *time.Time) []models.Flight {
	var filteredFlights []models.Flight

	if arrivalAirport != nil {
		for _, flight := range flights {
			flightArrival := *arrivalAirport
			if flight.Arrival == flightArrival {
				filteredFlights = append(filteredFlights, flight)
			}
		}
	}
	return filteredFlights
}
