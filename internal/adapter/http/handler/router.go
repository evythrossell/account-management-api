package handler

import (
	"github.com/evythrossell/account-management-api/internal/adapter/http/middleware"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func SetupRouter(
	accountHandler *AccountHandler,
	healthHandler *HealthHandler,
	transactionHandler *TransactionHandler,
) *gin.Engine {

	router := gin.Default()
	router.Use(middleware.Error())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
