package repositories

import (
	entities "flyhorizons-flightservice/repositories/entity"
	"flyhorizons-flightservice/services/interfaces"
)

type FlightRepository struct {
	*BaseRepository
}

var _ interfaces.FlightRepository = (*FlightRepository)(nil)

func NewFlightRepository(baseRepo *BaseRepository) *FlightRepository {
	return &FlightRepository{
		BaseRepository: baseRepo,
	}
}

func (repo *FlightRepository) GetAll() []entities.FlightEntity {
	db, _ := repo.CreateConnection()

	var flights []entities.FlightEntity
	db.Find(&flights)

	return flights
}

func (repo *FlightRepository) GetByFlightCode(flightCode string) entities.FlightEntity {
	db, _ := repo.CreateConnection()

	var flight entities.FlightEntity
	db.Where("FlightCode = ?", flightCode).First(&flight)

	return flight
}

func (repo *FlightRepository) Create(flightEntity entities.FlightEntity) entities.FlightEntity {
	db, _ := repo.CreateConnection()

	db.Create(&flightEntity)

	return flightEntity
}

func (repo *FlightRepository) DeleteByFlightCode(flightCode string) bool {
	db, _ := repo.CreateConnection()

	result := db.Where("FlightCode = ?", flightCode).Delete(&entities.FlightEntity{})

	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}

	return true
}

func (repo *FlightRepository) Update(flightEntity entities.FlightEntity) entities.FlightEntity {
	db, _ := repo.CreateConnection()

	db.Save(&flightEntity)

	return flightEntity
}
