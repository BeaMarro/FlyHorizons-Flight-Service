package mock_repositories

import (
	entities "flyhorizons-flightservice/repositories/entity"
	"flyhorizons-flightservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockFlightRepository struct {
	mock.Mock
}

var _ interfaces.FlightRepository = (*MockFlightRepository)(nil)

func (m *MockFlightRepository) GetByFlightCode(flightCode string) entities.FlightEntity {
	args := m.Called(flightCode)
	return args.Get(0).(entities.FlightEntity)
}

func (m *MockFlightRepository) GetAll() []entities.FlightEntity {
	args := m.Called()
	return args.Get(0).([]entities.FlightEntity)
}

func (m *MockFlightRepository) Create(flight entities.FlightEntity) entities.FlightEntity {
	args := m.Called(flight)
	return args.Get(0).(entities.FlightEntity)
}

func (m *MockFlightRepository) DeleteByFlightCode(flightCode string) bool {
	args := m.Called(flightCode)
	return args.Bool(0)
}

func (m *MockFlightRepository) Update(flight entities.FlightEntity) entities.FlightEntity {
	args := m.Called(flight)
	return args.Get(0).(entities.FlightEntity)
}
