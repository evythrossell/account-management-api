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
	"github.com/evythrossell/account-management-api/internal/adapter/http/middleware"
	"github.com/evythrossell/account-management-api/internal/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) GetAccount(ctx context.Context, documentNumber string) (*domain.Account, error) {
	args := m.Called(ctx, documentNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountService) GetAccountByDocument(ctx context.Context, documentNumber string) (*domain.Account, error) {
	args := m.Called(ctx, documentNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountService) CreateAccount(ctx context.Context, doc string) (*domain.Account, error) {
	args := m.Called(ctx, doc)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountService) GetAccountByID(ctx context.Context, id int64) (*domain.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func TestAccountHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("CreateAccount - Success", func(t *testing.T) {
		svc := new(MockAccountService)
		h := handler.NewAccountHandler(svc)
		r := gin.Default()
		r.POST("/accounts", h.CreateAccount)

		acc := &domain.Account{ID: 1, DocumentNumber: "123"}
		svc.On("CreateAccount", mock.Anything, "123").Return(acc, nil)

		body, _ := json.Marshal(handler.CreateAccountRequest{DocumentNumber: "123"})
		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("CreateAccount - Invalid JSON", func(t *testing.T) {
		h := handler.NewAccountHandler(nil)
		r := gin.Default()
		r.POST("/accounts", h.CreateAccount)

		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBufferString("{invalid}"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("CreateAccount - Service Error", func(t *testing.T) {
		svc := new(MockAccountService)
		h := handler.NewAccountHandler(svc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		svc.On("CreateAccount", mock.Anything, "123").Return(nil, errors.New("db error"))

		body, _ := json.Marshal(handler.CreateAccountRequest{DocumentNumber: "123"})
		c.Request, _ = http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))

		h.CreateAccount(c)

		assert.Len(t, c.Errors, 1)
		assert.Equal(t, "db error", c.Errors[0].Error())
		svc.AssertExpectations(t)
	})

	t.Run("GetAccount - Success", func(t *testing.T) {
		svc := new(MockAccountService)
		h := handler.NewAccountHandler(svc)
		r := gin.Default()
		r.GET("/accounts/:accountId", h.GetAccount)

		acc := &domain.Account{ID: 7, DocumentNumber: "123"}
		svc.On("GetAccountByID", mock.Anything, int64(7)).Return(acc, nil)

		req, _ := http.NewRequest("GET", "/accounts/7", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetAccount - Invalid ID", func(t *testing.T) {
		h := handler.NewAccountHandler(nil)
		r := gin.Default()
		r.Use(middleware.Error())
		r.GET("/accounts/:accountId", h.GetAccount)

		req, _ := http.NewRequest("GET", "/accounts/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetAccount - Service Error", func(t *testing.T) {
		svc := new(MockAccountService)
		h := handler.NewAccountHandler(svc)
		r := gin.Default()
		r.GET("/accounts/:accountId", h.GetAccount)

		svc.On("GetAccountByID", mock.Anything, int64(7)).Return(nil, errors.New("not found"))

		req, _ := http.NewRequest("GET", "/accounts/7", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}
