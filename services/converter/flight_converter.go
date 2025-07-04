package converter

import (
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/models/enums"
	entities "flyhorizons-flightservice/repositories/entity"
	"flyhorizons-flightservice/utils"
	"time"
)

type FlightConverter struct {
	WeekdayUtils utils.WeekdayUtils
}

func (flightConverter *FlightConverter) ConvertFlightEntityToFlight(entity entities.FlightEntity) models.Flight {
	// Convert JSON string to []enums.Day
	departureDays, err := flightConverter.WeekdayUtils.ConvertJSONToDays(entity.DepartureDays)
	if err != nil {
		departureDays = []enums.Day{}
	}

	return models.Flight{
		FlightCode:        entity.FlightCode,
		Departure:         entity.Departure,
		Arrival:           entity.Arrival,
		DurationInMinutes: entity.DurationInMinutes,
		DepartureTime:     entity.DepartureTime,
		DepartureDays:     departureDays,
		BasePrice:         entity.BasePrice,
	}
}

func (flightConverter *FlightConverter) ConvertFlightToFlightEntity(flight models.Flight) entities.FlightEntity {
	// Convert []enums.Day to JSON string
	departureDaysJSON, err := flightConverter.WeekdayUtils.ConvertDaysToJSON(flight.DepartureDays)
	if err != nil {
		departureDaysJSON = "[]"
	}

	return entities.FlightEntity{
		FlightCode:        flight.FlightCode,
		Departure:         flight.Departure,
		Arrival:           flight.Arrival,
		DurationInMinutes: flight.DurationInMinutes,
		DepartureTime:     flight.DepartureTime,
		DepartureDays:     departureDaysJSON,
		BasePrice:         flight.BasePrice,
		// Set current time for record creation/update
		CreatedAt: time.Now(),
	}
}
