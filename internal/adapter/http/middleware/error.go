package middleware

import (
	"errors"
	"net/http"

	common "github.com/evythrossell/account-management-api/pkg"
	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var de *common.DomainError
			if errors.As(err, &de) {
				c.AbortWithStatusJSON(de.HTTPStatusCode(), gin.H{
					"code":    de.Code,
					"message": de.PublicMessage(),
				})
				return
			}
			if errors.Is(err, common.ErrAccountNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"code":    "NOT_FOUND_ERROR",
					"message": "account not found",
				})
				return
			}
			if errors.Is(err, common.ErrTransactionNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"code":    "NOT_FOUND_ERROR",
					"message": "transaction not found",
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "an unexpected error occurred",
			})
		}
	}
}
