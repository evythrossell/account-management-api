package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/evythrossell/account-management-api/internal/adapter/http/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHealthService struct {
	mock.Mock
}

func (m *MockHealthService) Check(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHealthHandler_Check(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 200 when service is healthy", func(t *testing.T) {
		svc := new(MockHealthService)
		h := handler.NewHealthHandler(svc)

		svc.On("Check", mock.Anything).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/health", nil)

		h.Check(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "up")
		svc.AssertExpectations(t)
	})

	t.Run("should return 503 when service fails", func(t *testing.T) {
		svc := new(MockHealthService)
		h := handler.NewHealthHandler(svc)

		svc.On("Check", mock.Anything).Return(errors.New("db connection failed"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/health", nil)

		h.Check(c)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Contains(t, w.Body.String(), "unavailable")
		svc.AssertExpectations(t)
	})
}
