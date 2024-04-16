package httpserver

import (
	"cars-service/internal/app"
	"cars-service/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// New creates HTTP server with all needed routes
func New(addr string, a app.App, logs logger.Logger) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	setRoutes(api, a, logs)
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
