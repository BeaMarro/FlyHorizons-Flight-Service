package routes

import (
	"flyhorizons-flightservice/models"
	"flyhorizons-flightservice/services/errors"
	"flyhorizons-flightservice/services/interfaces"
	"flyhorizons-flightservice/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handles the flight CRUD functionality
func RegisterFlightRoutes(router *gin.Engine, flightService interfaces.FlightService, authMiddleware interfaces.GatewayAuthMiddleware) {
	// Public routes
	router.GET("/flights", func(ctx *gin.Context) {
		flights := flightService.GetAll(ctx.Request.Context())
		ctx.JSON(http.StatusOK, flights)
	})

	router.GET("flights/:flightCode", func(ctx *gin.Context) {
		flightCode := ctx.Param("flightCode")

		flight, err := flightService.GetByFlightCode(ctx.Request.Context(), flightCode)

		if err != nil {
			if _, ok := err.(*errors.FlightNotFoundError); ok {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()}) // 404 Not Found
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, flight)
	})

	flightGroup := router.Group("/flights")
	flightGroup.Use(authMiddleware.GatewayAuthMiddleware())

	// Protected routes
	// Only accessible by admins
	flightGroup.POST("/", utils.IPWhitelistingMiddleware(), func(ctx *gin.Context) { // Whitelists IP addresses to only accept the admin ones
		role, exists := ctx.Get("role")
		fmt.Println("Loggedin role: ", role)

		if !exists || role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: admin access required"})
			return
		}

		var flight models.Flight
		if err := ctx.ShouldBindJSON(&flight); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		postFlight, err := flightService.Create(ctx.Request.Context(), flight)
		if err != nil {
			if _, ok := err.(*errors.FlightExistsError); ok {
				ctx.JSON(http.StatusConflict, gin.H{"message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, postFlight)
	})

	flightGroup.DELETE("/:flightCode", func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		fmt.Println("Loggedin role: ", role)

		if !exists || role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: admin access required"})
			return
		}

		flightCode := ctx.Param("flightCode")

		success, err := flightService.DeleteByFlightCode(ctx.Request.Context(), flightCode)
		if err != nil {
			if _, ok := err.(*errors.FlightNotFoundError); ok {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()}) // 404 Not Found
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		// Uses success to confirm the deletion
		if success {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Flight deleted successfully",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to delete flight, but no error has occurred",
			})
		}
	})

	flightGroup.PUT("/", func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		fmt.Println("Loggedin role: ", role)

		if !exists || role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: admin access required"})
			return
		}

		var flight models.Flight
		// Convert the JSON to a Flight object
		if err := ctx.ShouldBindJSON(&flight); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		put_flight, err := flightService.Update(ctx.Request.Context(), flight)
		if err != nil {
			if _, ok := err.(*errors.FlightNotFoundError); ok {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()}) // 404 Not Found
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, put_flight)
	})
}
