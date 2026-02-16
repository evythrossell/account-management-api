package http

import (
	"errors"
	"net/http"

	domainerror "github.com/evythrossell/account-management-api/internal/core/error"
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

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request format",
		})
		return
	}

	account, err := h.service.CreateAccount(c.Request.Context(), req.DocumentNumber)
	if err != nil {
		var de *domainerror.DomainError
		if errors.As(err, &de) {
			c.JSON(de.HTTPStatusCode(), ErrorResponse{
				Code:    de.Code,
				Message: de.PublicMessage(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "An unexpected error occurred",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, account)
}
