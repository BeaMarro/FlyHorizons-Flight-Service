package interfaces

import (
	entities "flyhorizons-flightservice/repositories/entity"
)

type FlightRepository interface {
	GetAll() []entities.FlightEntity
	GetByFlightCode(flightCode string) entities.FlightEntity
	Create(flight entities.FlightEntity) entities.FlightEntity
	DeleteByFlightCode(flightCode string) bool
	Update(flight entities.FlightEntity) entities.FlightEntity
}
