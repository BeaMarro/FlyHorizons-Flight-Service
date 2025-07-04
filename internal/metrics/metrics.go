package metrics

import (
	"flyhorizons-flightservice/internal/health"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetricsRoutes(router *gin.Engine, dbCheck health.DatabaseCheck) {
	dbHealthGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mssql_db_health",
		Help: "Database health status: 1 for up, 0 for down",
	})

	prometheus.MustRegister(dbHealthGauge)

	go func() {
		for {
			if dbCheck.Pass() {
				dbHealthGauge.Set(1)
			} else {
				dbHealthGauge.Set(0)
			}
			// sleep some interval before next check
			// this is configured to check it every 10 seconds
			time.Sleep(10 * time.Second)
		}
	}()

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
