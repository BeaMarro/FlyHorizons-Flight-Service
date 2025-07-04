package strategies

import (
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/models/enums"
	"flyhorizons-flightservice/utils"
	"time"
)

type DateRangeStrategy struct {
}

func (strategy DateRangeStrategy) Filter(flights []models.Flight, depatureAirport *string, arrivalAirport *string, departureDate *time.Time, returnDate *time.Time) []models.Flight {
	var filteredFlights []models.Flight
	weekdayUtils := utils.WeekdayUtils{}

	if departureDate != nil {
		// Convert departureDate to Day enum
		flightDepatureDate := *departureDate
		departureWeekday := weekdayUtils.ConvertToWeekDay(flightDepatureDate)

		// Convert arrivalDate to Day enum (if provided)
		var arrivalWeekday *enums.Day

		if returnDate != nil {
			tempArrivalDay := weekdayUtils.ConvertToWeekDay(*returnDate)
			arrivalWeekday = &tempArrivalDay
		}

		for _, flight := range flights {
			// Check if the flight departure day matches any of the flight allowed days
			if weekdayUtils.ContainsDay(flight.DepartureDays, departureWeekday) {
				// If arrivalDate is provided, ensure that the flight departure day matches with the arrival weekday
				if arrivalWeekday == nil || departureWeekday == *arrivalWeekday {
					filteredFlights = append(filteredFlights, flight)
				}
			}
		}
	}

	return filteredFlights
}
