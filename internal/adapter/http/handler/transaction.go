package handler

import (
	"net/http"
	"strconv"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	"github.com/evythrossell/account-management-api/internal/core/port"
	"github.com/gin-gonic/gin"
)

type createTransactionRequest struct {
	AccountID     int64   `json:"account_id" binding:"required" example:"123"`
	OperationType int16   `json:"operation_type_id" binding:"required" example:"1"`
	Amount        float64 `json:"amount" binding:"required" example:"100.50"`
}

// Tipos de erro específicos para cada status code
// BadRequestError, NotFoundError, InternalServerError estão definidos em account.go e são reutilizados aqui

type TransactionHandler struct {
	service port.TransactionService
}

func NewTransactionHandler(service port.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

// CreateTransaction godoc
// @Summary      Criar transação
// @Description  Cria uma nova transação bancária (débito/crédito)
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        body body createTransactionRequest true "Dados da transação"
// @Success      201 {object} domain.Transaction "Transação criada com sucesso"
// @Failure      400 {object} BadRequestError "Erro de validação"
// @Failure      404 {object} NotFoundError "Conta não encontrada"
// @Failure      500 {object} InternalServerError "Erro interno do servidor"
// @Router       /v1/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    domain.ErrCodeInvalidBody,
			"message": domain.ErrMsgInvalidBodyRequest,
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

// GetTransaction godoc
// @Summary      Obter transação por ID
// @Description  Retorna os detalhes de uma transação específica
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        transactionId path int64 true "ID da transação"
// @Success      200 {object} domain.Transaction "Transação encontrada"
// @Failure      400 {object} BadRequestError "ID inválido"
// @Failure      404 {object} NotFoundError "Transação não encontrada"
// @Failure      500 {object} InternalServerError "Erro interno do servidor"
// @Router       /v1/transactions/{transactionId} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("transactionId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    domain.ErrCodeInvalidID,
			"message": domain.ErrMsgTransactionIDInvalid,
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
