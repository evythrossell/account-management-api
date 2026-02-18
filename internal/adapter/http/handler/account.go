package handler

import (
	"net/http"
	"strconv"

	"github.com/evythrossell/account-management-api/internal/core/port"
	"github.com/gin-gonic/gin"

	common "github.com/evythrossell/account-management-api/pkg"
)

type AccountHandler struct {
	service port.AccountService
}

func NewAccountHandler(service port.AccountService) *AccountHandler {
	return &AccountHandler{
		service: service,
	}
}

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" binding:"required"`
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": "INVALID_BODY", "message": "invalid request body or missing required fields"})
		return
	}

	account, err := h.service.CreateAccount(c.Request.Context(), req.DocumentNumber)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	idParam := c.Param("accountId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Error(common.NewValidationError("the account ID must be a valid integer", err))
		return
	}

	account, err := h.service.GetAccountByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}
