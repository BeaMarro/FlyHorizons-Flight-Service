package main

import (
	cache "flyhorizons-flightservice/config"
	"flyhorizons-flightservice/internal/health"
	"flyhorizons-flightservice/internal/metrics"
	"flyhorizons-flightservice/utils"

	"flyhorizons-flightservice/repositories"
	"flyhorizons-flightservice/routes"
	"flyhorizons-flightservice/services"
	"flyhorizons-flightservice/services/authentication"
	"flyhorizons-flightservice/services/converter"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	router := gin.Default()
	baseRepo := repositories.BaseRepository{}
	dbCheck := health.DatabaseCheck{Repository: &baseRepo}
	_ = godotenv.Load()

	// Health check setup
	conf := config.DefaultConfig()
	healthcheck.New(router, conf, []checks.Check{dbCheck})

	// Metrics setup
	metrics.RegisterMetricsRoutes(router, dbCheck)

	// Microservice setup
	redis := cache.CreateRedisClient()
	utils.LoadWhitelistedIPs()
	flightRepo := repositories.NewFlightRepository(&baseRepo)
	flightConverter := converter.FlightConverter{}

	gatewayAuthMiddleware := authentication.NewGatewayAuthMiddleware()
	flightService := services.NewFlightService(flightRepo, flightConverter, redis)

	routes.RegisterFlightRoutes(router, flightService, gatewayAuthMiddleware)
	routes.RegisterFilterFlightRoutes(router, flightService)

	router.Run(":8080")
}
