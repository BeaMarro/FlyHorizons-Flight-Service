package errors

import "fmt"

type FlightExistsError struct {
	FlightCode string
}

func (e *FlightExistsError) Error() string {
	return fmt.Sprintf("Flight with the code %s already exists", e.FlightCode)
}

func NewFlightExistsError(flightCode string, errorCode int) *FlightExistsError {
	return &FlightExistsError{FlightCode: flightCode}
}
