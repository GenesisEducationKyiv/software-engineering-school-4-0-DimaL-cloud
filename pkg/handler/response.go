package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type error struct {
	Message string `json:"message"`
}

func newError(c *gin.Context, statusCode int, message string) {
	log.Error(message)
	c.AbortWithStatusJSON(statusCode, error{Message: message})
}
