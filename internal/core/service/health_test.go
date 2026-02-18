package services_test

import (
	"context"
	"errors"
	"testing"

	services "github.com/evythrossell/account-management-api/internal/core/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDBHealthChecker struct {
	mock.Mock
}

func (m *MockDBHealthChecker) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHealthService(t *testing.T) {
	ctx := context.Background()

	t.Run("Check - Success", func(t *testing.T) {
		mockChecker := new(MockDBHealthChecker)
		svc := services.NewHealthService(mockChecker)

		mockChecker.On("PingContext", ctx).Return(nil)

		err := svc.Check(ctx)

		assert.NoError(t, err)
		mockChecker.AssertExpectations(t)
	})

	t.Run("Check - Failure", func(t *testing.T) {
		mockChecker := new(MockDBHealthChecker)
		svc := services.NewHealthService(mockChecker)

		mockChecker.On("PingContext", ctx).Return(errors.New("db connection failed"))

		err := svc.Check(ctx)

		assert.Error(t, err)
		assert.Equal(t, "db connection failed", err.Error())
		mockChecker.AssertExpectations(t)
	})
}
