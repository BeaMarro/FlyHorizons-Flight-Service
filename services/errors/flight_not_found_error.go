package errors

import "fmt"

type FlightNotFoundError struct {
	FlightCode string
}

func (e *FlightNotFoundError) Error() string {
	return fmt.Sprintf("Flight with the code %s was not found", e.FlightCode)
}

func NewFlightNotFoundError(flightCode string, errorCode int) *FlightNotFoundError {
	return &FlightNotFoundError{FlightCode: flightCode}
}
