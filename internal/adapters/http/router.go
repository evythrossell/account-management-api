package http

import (
	"github.com/evythrossell/account-management-api/internal/adapters/http/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	accountHandler *AccountHandler,
	healthHandler *HealthHandler,
	transactionHandler *TransactionHandler,
) *gin.Engine {

	router := gin.Default()
	router.Use(middleware.Error())

	router.GET("/health", healthHandler.Check)

	v1 := router.Group("/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("/:accountId", accountHandler.GetAccount)
		}

		transactions := v1.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("/:transactionId", transactionHandler.GetTransaction)
		}
	}

	return router
}
