package httpserver

import (
	_ "cars-service/docs"
	"cars-service/internal/app"
	"cars-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func setRoutes(r *gin.RouterGroup, a app.App, logs logger.Logger) {
	r.Use(panicMiddleware(logs))
	r.Use(loggingMiddleware(logs))

	r.GET("/cars/:id", handleGetCarById(a))
	r.GET("/cars", handleGetCars(a))
	r.POST("/cars", handleAddCars(a))
	r.PUT("/cars/:id", handleUpdateCar(a))
	r.DELETE("/cars/:id", handleDeleteCar(a))
}
