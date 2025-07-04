package mock_repositories

import (
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockFlightService struct {
	mock.Mock
}

var _ interfaces.FlightService = (*MockFlightService)(nil)

func (m *MockFlightService) GetAll() []models.Flight {
	args := m.Called()
	return args.Get(0).([]models.Flight)
}

func (m *MockFlightService) GetByFlightCode(flightCode string) (*models.Flight, error) {
	args := m.Called(flightCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Flight), args.Error(1)
}

func (m *MockFlightService) FlightExists(flightCode string) bool {
	args := m.Called(flightCode)
	return args.Bool(0)
}

func (m *MockFlightService) Create(flight models.Flight) (*models.Flight, error) {
	args := m.Called(flight)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Flight), args.Error(1)
}

func (m *MockFlightService) DeleteByFlightCode(flightCode string) (bool, error) {
	args := m.Called(flightCode)
	return args.Bool(0), args.Error(1)
}

func (m *MockFlightService) Update(user models.Flight) (*models.Flight, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Flight), args.Error(1)
}
