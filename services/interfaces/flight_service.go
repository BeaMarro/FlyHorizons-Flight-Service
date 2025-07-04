package interfaces

import (
	"flyhorizons-flightservice/models"

	"golang.org/x/net/context"
)

type FlightService interface {
	GetAll(ctx context.Context) []models.Flight
	GetByFlightCode(ctx context.Context, flightCode string) (*models.Flight, error)
	FlightExists(ctx context.Context, flightCode string) bool
	Create(ctx context.Context, flight models.Flight) (*models.Flight, error)
	DeleteByFlightCode(ctx context.Context, flightCode string) (bool, error)
	Update(ctx context.Context, flight models.Flight) (*models.Flight, error)
}
