package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/evythrossell/account-management-api/internal/adapter/http/middleware"
	common "github.com/evythrossell/account-management-api/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should handle domain error - Not Found", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.Error())

		domainErr := common.NewNotFoundError("account not found", nil)

		r.GET("/not-found", func(c *gin.Context) {
			c.Error(domainErr)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/not-found", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "NOT_FOUND_ERROR")
		assert.Contains(t, w.Body.String(), "account not found")
	})

	t.Run("should handle domain error - Validation", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.Error())

		domainErr := common.NewValidationError("invalid input", nil)

		r.GET("/validation-error", func(c *gin.Context) {
			c.Error(domainErr)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/validation-error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "VALIDATION_ERROR")
	})

	t.Run("should handle generic error as 500", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.Error())

		r.GET("/generic-error", func(c *gin.Context) {
			c.Error(errors.New("unexpected database failure"))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/generic-error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "INTERNAL_SERVER_ERROR")
		assert.Contains(t, w.Body.String(), "an unexpected error occurred")
	})

	t.Run("should do nothing when there are no errors", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.Error())

		r.GET("/success", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/success", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Body.String())
	})

	t.Run("should handle ErrAccountNotFound specifically", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.Error())

		r.GET("/account-not-found", func(c *gin.Context) {
			c.Error(common.ErrAccountNotFound)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/account-not-found", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "NOT_FOUND_ERROR")
		assert.Contains(t, w.Body.String(), "account not found")
	})

	t.Run("should handle ErrTransactionNotFound specifically", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.Error())

		r.GET("/transaction-not-found", func(c *gin.Context) {
			c.Error(common.ErrTransactionNotFound)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/transaction-not-found", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "NOT_FOUND_ERROR")
		assert.Contains(t, w.Body.String(), "transaction not found")
	})

	t.Run("should handle ErrInvalidAmount specifically", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.Error())

		r.GET("/invalid-amount", func(c *gin.Context) {
			c.Error(common.ErrInvalidAmount)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/invalid-amount", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "VALIDATION_ERROR")
		assert.Contains(t, w.Body.String(), "amount must be greater than zero")
	})

	t.Run("should handle ErrInvalidOperation specifically", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.Error())

		r.GET("/invalid-operation", func(c *gin.Context) {
			c.Error(common.ErrInvalidOperation)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/invalid-operation", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "VALIDATION_ERROR")
		assert.Contains(t, w.Body.String(), "invalid operation type")
	})
}
