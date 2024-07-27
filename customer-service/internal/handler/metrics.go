package handler

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
)

func (h *Handler) metrics(c *gin.Context) {
	metrics.WritePrometheus(c.Writer, true)
}
