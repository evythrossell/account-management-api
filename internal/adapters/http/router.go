package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(accountHandler *AccountHandler, healthHandler *HealthHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/accounts", accountHandler.CreateAccount)
	router.GET("/accounts/:accountId", accountHandler.GetAccount)

	router.GET("/health", healthHandler.Check)

	return router
}
