package http

import (
	"context"
	"errors"
	"net/http"
	"time"
	"strconv"

	domainerror "github.com/evythrossell/account-management-api/internal/core/error"
	"github.com/evythrossell/account-management-api/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type createTransactionRequest struct {
	AccountID     int64   `json:"account_id"`
	OperationType int16   `json:"operation_type_id"`
	Amount        float64 `json:"amount"`
}

type TransactionHandler struct {
	t ports.TransactionService
}

func NewTransactionHandler(t ports.TransactionService) *TransactionHandler {
	return &TransactionHandler{t: t}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "invalid request body"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	tx, err := h.t.CreateTransaction(ctx, req.AccountID, req.OperationType, req.Amount)
	if err != nil {
		var de *domainerror.DomainError
		if errors.As(err, &de) {
			c.JSON(de.HTTPStatusCode(), gin.H{"code": de.Code, "message": de.PublicMessage()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, tx)
}


func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	transactionID, err := strconv.ParseInt(c.Param("transactionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_ERROR", "message": "invalid transactionId"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	transaction, err := h.t.GetByTransactionID(
		ctx,
		transactionID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, transaction)
}