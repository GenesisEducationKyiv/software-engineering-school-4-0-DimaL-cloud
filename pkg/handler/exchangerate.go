package handler

import (
	"exchange-rate-notifier-api/pkg/client"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) rate(c *gin.Context) {
	rateClient := client.NewExchangeRateClient()
	rate, err := rateClient.GetCurrentExchangeRate()
	if err != nil {
		newError(c, http.StatusBadRequest, "Invalid status value")
	}
	c.String(http.StatusOK, fmt.Sprintf("%f", rate.Rate))
}
