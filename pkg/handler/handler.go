package handler

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		api.POST("/subscribe", h.subscribe)
		api.POST("/unsubscribe", h.unsubscribe)
		api.GET("/rate", h.rate)
	}
	return router
}
