package handler

import (
	"net/http"
	"strconv"

	"github.com/evythrossell/account-management-api/internal/core/domain"
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
	DocumentNumber string `json:"document_number" binding:"required" example:"12345678901"`
}

type BadRequestError struct {
	Code    string `json:"code" example:"VALIDATION_ERROR"`
	Message string `json:"message" example:"document_number is required"`
}

type NotFoundError struct {
	Code    string `json:"code" example:"NOT_FOUND"`
	Message string `json:"message" example:"account not found"`
}

type InternalServerError struct {
	Code    string `json:"code" example:"INTERNAL_ERROR"`
	Message string `json:"message" example:"unexpected error on internal service"`
}

// CreateAccount godoc
// @Summary      Criar nova conta
// @Description  Cria uma nova conta bancária com CPF ou CNPJ
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        body body CreateAccountRequest true "Dados da conta"
// @Success      201 {object} domain.Account "Conta criada com sucesso"
// @Failure      400 {object} BadRequestError "Erro de validação"
// @Failure      500 {object} InternalServerError "Erro interno do servidor"
// @Router       /v1/accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": domain.ErrCodeInvalidBody, "message": domain.ErrMsgInvalidBodyRequest})
		return
	}

	account, err := h.service.CreateAccount(c.Request.Context(), req.DocumentNumber)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccount godoc
// @Summary      Obter conta por ID
// @Description  Retorna os detalhes de uma conta específica
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        accountId path int64 true "ID da conta"
// @Success      200 {object} domain.Account "Conta encontrada"
// @Failure      400 {object} BadRequestError "ID inválido"
// @Failure      404 {object} NotFoundError "Conta não encontrada"
// @Failure      500 {object} InternalServerError "Erro interno do servidor"
// @Router       /v1/accounts/{accountId} [get]
func (h *AccountHandler) GetAccount(c *gin.Context) {
	idParam := c.Param("accountId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Error(common.NewValidationError(domain.ErrMsgAccountIDInvalid, err))
		return
	}

	account, err := h.service.GetAccountByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}
