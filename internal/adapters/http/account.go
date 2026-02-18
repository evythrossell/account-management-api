package http

import (
	"net/http"
	"strconv"

	"github.com/evythrossell/account-management-api/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	service ports.AccountService
}

func NewAccountHandler(service ports.AccountService) *AccountHandler {
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": "INVALID_BODY", "message": err.Error()})
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
	id, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": "INVALID_ID", "message": "the account ID must be an integer"})
		return
	}

	account, err := h.service.GetAccountByID(c.Request.Context(), int64(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}
