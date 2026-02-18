package handler

import (
	"net/http"

	"github.com/evythrossell/account-management-api/internal/core/port"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	service port.HealthService
}

func NewHealthHandler(service port.HealthService) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	if err := h.service.Check(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unavailable",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
