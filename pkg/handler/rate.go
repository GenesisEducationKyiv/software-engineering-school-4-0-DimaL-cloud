package handler

import (
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/client"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Get current exchange rate
// @Description Get the current exchange rate from the external API
// @Tags rate
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Returns the current exchange rate"
// @Failure 500 {object} handler.errorResponse "failed to get exchange rate"
// @Router /rate [get]
func (h *Handler) rate(c *gin.Context) {
	rateClient := client.NewExchangeRateClient()
	rate, err := rateClient.GetCurrentExchangeRate()
	if err != nil {
		newError(c, http.StatusInternalServerError, "failed to get exchange rate")
	}
	c.String(http.StatusOK, fmt.Sprintf("%f", rate.Rate))
}
