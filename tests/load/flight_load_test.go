package load_test

import (
	load_test_utils "flyhorizons-flightservice/tests/load/utils"
	"fmt"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// Before running the load tests
// Run the microservice at the same time

type FlightLoadTest struct {
	loadTestUtils load_test_utils.LoadTestUtils
}

func getAllFlights(t *testing.T, rate vegeta.Rate, duration time.Duration, htmlReport string, title string) vegeta.Metrics {
	loadTest := FlightLoadTest{
		loadTestUtils: load_test_utils.LoadTestUtils{},
	}

	target := vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/flights",
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(vegeta.NewStaticTargeter(target), rate, duration, "Load Test GetAllFlights") {
		metrics.Add(res)
	}
	metrics.Close()

	// Log the metrics
	loadTest.loadTestUtils.LogMetrics(t, &metrics)
	loadTest.loadTestUtils.GenerateHTMLReport(t, &metrics, fmt.Sprintf("%s.html", htmlReport), fmt.Sprintf("Load Test: %s", title))

	return metrics
}

// DONE: Load Test (short time): Get all flights at 10 requests per second for 1 second
func TestLoadGetAllFlightsFewTimesUsingShortTime(t *testing.T) {
	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 1 * time.Second
	htmlReport := "get_few_flights_short_time"
	title := "Get all flights at 10 requests per second for 1 second"
	_ = getAllFlights(t, rate, duration, htmlReport, title)
}

// TODO: Load Test (short time): Get all flights at 200 requests per second for 1 second
func TestLoadGetAllFlightsManyTimesUsingShortTime(t *testing.T) {
	rate := vegeta.Rate{Freq: 200, Per: time.Second}
	duration := 1 * time.Second
	htmlReport := "get_many_flights_short_time"
	title := "Get all flights at 200 requests per second for 1 second"
	_ = getAllFlights(t, rate, duration, htmlReport, title)
}

// DONE: Load Test (long time): Get all flights at 10 requests per second for 10 seconds
func TestLoadGetAllFlightsFewTimesUsingLongerTime(t *testing.T) {
	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 10 * time.Second
	htmlReport := "get_many_flights_long_time"
	title := "Get all flights at 10 requests per second for 10 seconds"
	_ = getAllFlights(t, rate, duration, htmlReport, title)
}

// TODO: Spike Test: Get all flights at 500 requests per second for 1 second
func TestLoadSpikeGetAllFlightsVeryHigh(t *testing.T) {
	rate := vegeta.Rate{Freq: 500, Per: time.Second}
	duration := 1 * time.Second
	htmlReport := "get_flights_spike"
	title := "Spike test get all flights at 500 requests per second for 1 second"
	_ = getAllFlights(t, rate, duration, htmlReport, title)
}

// TODO: Spike Test: Get all flights at 1000 requests per second for 1 second
func TestLoadSpikeGetAllFlightsExtremelyHigh(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 1 * time.Second
	htmlReport := "get_flights_spike_extreme"
	title := "Spike test get all flights at 1000 requests per second for 1 second"
	_ = getAllFlights(t, rate, duration, htmlReport, title)
}
