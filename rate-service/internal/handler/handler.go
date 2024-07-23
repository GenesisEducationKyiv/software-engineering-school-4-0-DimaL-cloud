package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"rate-service/internal/service"
)

type Handler struct {
	rateService         service.Rate
	subscriptionService service.Subscription
}

func NewHandler(rateService service.Rate, subscriptionService service.Subscription) *Handler {
	return &Handler{
		rateService:         rateService,
		subscriptionService: subscriptionService,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/metrics", h.metrics)
	api := router.Group("/api")
	{
		api.POST("/subscribe", h.subscribe)
		api.POST("/unsubscribe", h.unsubscribe)
		api.GET("/rate", h.rate)
	}
	return router
}
