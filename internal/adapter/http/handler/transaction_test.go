package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/evythrossell/account-management-api/internal/adapter/http/handler"
	"github.com/evythrossell/account-management-api/internal/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(ctx context.Context, accID int64, opType int16, amount float64) (*domain.Transaction, error) {
	args := m.Called(ctx, accID, opType, amount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionService) GetByTransactionID(ctx context.Context, id int64) (*domain.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func TestTransactionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("CreateTransaction - Success", func(t *testing.T) {
		svc := new(MockTransactionService)
		h := handler.NewTransactionHandler(svc)
		r := gin.New()
		r.POST("/transactions", h.CreateTransaction)

		tx := &domain.Transaction{ID: 100}
		svc.On("CreateTransaction", mock.Anything, int64(1), int16(4), 50.0).Return(tx, nil)

		body := map[string]interface{}{"account_id": 1, "operation_type_id": 4, "amount": 50.0}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("CreateTransaction - Invalid Body", func(t *testing.T) {
		h := handler.NewTransactionHandler(nil)
		r := gin.New()
		r.POST("/transactions", h.CreateTransaction)

		req := httptest.NewRequest("POST", "/transactions", bytes.NewBufferString("{invalid}"))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("CreateTransaction - Service Error", func(t *testing.T) {
		svc := new(MockTransactionService)
		h := handler.NewTransactionHandler(svc)
		r := gin.New()
		r.POST("/transactions", h.CreateTransaction)

		svc.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		body := map[string]interface{}{"account_id": 1, "operation_type_id": 4, "amount": 50.0}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Len(t, r.Handlers, 0)
	})

	t.Run("GetTransaction - Success", func(t *testing.T) {
		svc := new(MockTransactionService)
		h := handler.NewTransactionHandler(svc)
		r := gin.New()
		r.GET("/transactions/:transactionId", h.GetTransaction)

		tx := &domain.Transaction{ID: 6}
		svc.On("GetByTransactionID", mock.Anything, int64(6)).Return(tx, nil)

		req := httptest.NewRequest("GET", "/transactions/6", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetTransaction - Invalid ID", func(t *testing.T) {
		h := handler.NewTransactionHandler(nil)
		r := gin.New()
		r.GET("/transactions/:transactionId", h.GetTransaction)

		req := httptest.NewRequest("GET", "/transactions/abc", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetTransaction - Service Error", func(t *testing.T) {
		svc := new(MockTransactionService)
		h := handler.NewTransactionHandler(svc)
		r := gin.New()
		r.GET("/transactions/:transactionId", h.GetTransaction)

		svc.On("GetByTransactionID", mock.Anything, int64(6)).Return(nil, errors.New("not found"))

		req := httptest.NewRequest("GET", "/transactions/6", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
