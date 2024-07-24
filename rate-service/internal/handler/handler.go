package handler

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
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
	router.Use(HTTPRequestMetricsMiddleware())
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

func HTTPRequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		c.Next()
		statusCode := c.Writer.Status()
		metrics.GetOrCreateCounter(fmt.Sprintf(`http_requests_total{path="%s", status="%d"}`, path, statusCode)).Inc()
	}
}
