package entities

import (
	"time"
)

type FlightEntity struct {
	FlightCode        string    `gorm:"column:FlightCode;primaryKey"`
	Departure         string    `gorm:"column:Departure"`
	Arrival           string    `gorm:"column:Arrival"`
	DurationInMinutes int       `gorm:"column:DurationInMinutes"`
	DepartureTime     time.Time `gorm:"column:DepartureTime"`
	DepartureDays     string    `gorm:"column:DepartureDays;type:string"` // JSON list of integers (string)
	BasePrice         float32   `gorm:"column:BasePrice"`
	CreatedAt         time.Time `gorm:"column:CreatedAt"`
}

// Override the default table name
func (FlightEntity) TableName() string {
	return "Flight"
}
