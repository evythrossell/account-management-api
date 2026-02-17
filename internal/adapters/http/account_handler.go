package http

import (
	"errors"
	"net/http"
	"strconv"

	domainerror "github.com/evythrossell/account-management-api/internal/core/error"
	"github.com/evythrossell/account-management-api/internal/core/ports"
	validator "github.com/evythrossell/account-management-api/internal/core/domain/validator"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if !validator.IsValidDocument(req.DocumentNumber) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "document_number must contain only digits and be 11 or 14 chars long"
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

func (h *AccountHandler) GetAccount(c *gin.Context) {
	accountIdStr := c.Param("accountId")

	if accountIdStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "accountId is required",
		})
		return
	}

	accountID, parseErr := strconv.ParseInt(accountIdStr, 10, 64)
	if parseErr != nil || accountID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "invalid accountId",
		})
		return
	}

	account, err := h.service.GetAccountByID(c.Request.Context(), accountID)
	if err == nil {
		c.JSON(http.StatusOK, account)
		return
	}

	var de *domainerror.DomainError
	if errors.As(err, &de) {
		c.JSON(de.HTTPStatusCode(), ErrorResponse{
			Code:    de.Code,
			Message: de.PublicMessage(),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "An unexpected error occurred",
	})
}
