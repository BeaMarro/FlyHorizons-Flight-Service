package models

import (
	"flyhorizons-flightservice/models/enums"
	"time"
)

type Flight struct {
	FlightCode        string      `json:"flight_code"`
	Departure         string      `json:"departure"`
	Arrival           string      `json:"arrival"`
	DurationInMinutes int         `json:"duration_in_minutes"`
	DepartureTime     time.Time   `json:"departure_time"`
	DepartureDays     []enums.Day `json:"departure_days"`
	BasePrice         float32     `json:"base_price"`
}
