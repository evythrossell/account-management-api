package http

import (
	"net/http"

	"github.com/evythrossell/account-management-api/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	service ports.AccountService
}

func NewAccountHandler(s ports.AccountService) *AccountHandler {
	return &AccountHandler{service: s}
}

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" binding:"required"`
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format."})
		return
	}

	account, err := h.service.CreateAccount(c.Request.Context(), req.DocumentNumber)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}
