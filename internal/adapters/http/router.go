package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(accountHandler *AccountHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/accounts", accountHandler.CreateAccount)
	router.GET("/accounts/:accountId", accountHandler.GetAccount)

	return router
}
