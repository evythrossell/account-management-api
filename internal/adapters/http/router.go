package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(accountHandler *AccountHandler, healthHandler *HealthHandler, transactionHandler *TransactionHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/accounts", accountHandler.CreateAccount)
	router.GET("/accounts/:accountId", accountHandler.GetAccount)

	router.POST("/transactions", transactionHandler.CreateTransaction)
	router.GET("/transactions/:transactionId", transactionHandler.GetTransaction)

	router.GET("/health", healthHandler.Check)

	return router
}
