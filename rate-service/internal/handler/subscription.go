package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
	"rate-service/internal/models"
)

// @Summary Subscribe to notifications
// @Tags subscription
// @Description subscribe to notifications
// @Param email query string true "email"
// @Success 200 "ok"
// @Failure 400 {object} handler.errorResponse "email is empty"
// @Failure 400 {object} handler.errorResponse "invalid email format"
// @Failure 409 {object} handler.errorResponse "subscription already exists"
// @Failure 500 {object} handler.errorResponse "failed to create subscription"
// @Router /api/subscribe [post]
func (h *Handler) subscribe(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		newError(c, http.StatusBadRequest, "email is empty")
		return
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		newError(c, http.StatusBadRequest, "invalid email format")
		return
	}
	err = h.subscriptionService.CreateSubscription(email)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			newError(c, http.StatusConflict, "subscription already exists")
			return
		}
		newError(c, http.StatusInternalServerError, "failed to create subscription")
		return
	}
}

// @Summary Unsubscribe from notifications
// @Tags subscription
// @Description unsubscribe from notifications
// @Param email query string true "email"
// @Success 200 "ok"
// @Failure 400 {string} string "email is empty"
// @Failure 400 {string} string "invalid email format"
// @Failure 500 {string} string "failed to delete subscription"
// @Router /api/unsubscribe [post]
func (h *Handler) unsubscribe(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		newError(c, http.StatusBadRequest, "email is empty")
		return
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		newError(c, http.StatusBadRequest, "invalid email format")
		return
	}
	err = h.subscriptionService.DeleteSubscription(email)
	if err != nil {
		newError(c, http.StatusInternalServerError, "failed to delete subscription")
		return
	}
}
