package handler

import (
	"net/http"

	"github.com/evythrossell/account-management-api/internal/core/port"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	service port.HealthService
}

type HealthResponse struct {
	Status string `json:"status" example:"up"`
	Detail string `json:"detail,omitempty" example:"database connection failed"`
}

type ServiceUnavailableError struct {
	Status string `json:"status" example:"unavailable"`
	Detail string `json:"detail" example:"database connection failed"`
}

func NewHealthHandler(service port.HealthService) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

// Check godoc
// @Summary      Verificar saúde da API
// @Description  Retorna o status de saúde da API e suas dependências
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200 {object} HealthResponse "API está operacional"
// @Failure      503 {object} ServiceUnavailableError "API indisponível"
// @Router       /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	if err := h.service.Check(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unavailable",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "up"})
}
