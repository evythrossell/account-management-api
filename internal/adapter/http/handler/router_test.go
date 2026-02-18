package handler_test

import (
	"testing"

	"github.com/evythrossell/account-management-api/internal/adapter/http/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should initialize router with all routes", func(t *testing.T) {
		accHandler := handler.NewAccountHandler(nil)
		healthHandler := handler.NewHealthHandler(nil)
		transHandler := handler.NewTransactionHandler(nil)

		r := handler.SetupRouter(accHandler, healthHandler, transHandler)

		assert.NotNil(t, r)

		routes := r.Routes()
		expectedRoutes := []string{
			"/health",
			"/v1/accounts",
			"/v1/accounts/:accountId",
			"/v1/transactions",
			"/v1/transactions/:transactionId",
		}

		for _, expected := range expectedRoutes {
			found := false
			for _, route := range routes {
				if route.Path == expected {
					found = true
					break
				}
			}
			assert.True(t, found, "route %s not found", expected)
		}
	})
}
