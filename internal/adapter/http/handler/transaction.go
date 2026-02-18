package handler

import (
	"net/http"
	"strconv"

	"github.com/evythrossell/account-management-api/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type createTransactionRequest struct {
	AccountID     int64   `json:"account_id" binding:"required"`
	OperationType int16   `json:"operation_type_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}

type TransactionHandler struct {
	service ports.TransactionService
}

func NewTransactionHandler(service ports.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_BODY",
			"message": "invalid request body or missing required fields",
		})
		return
	}

	tx, err := h.service.CreateTransaction(c.Request.Context(), req.AccountID, req.OperationType, req.Amount)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("transactionId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "the transaction ID must be a valid integer",
		})
		return
	}

	transaction, err := h.service.GetByTransactionID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, transaction)
}
