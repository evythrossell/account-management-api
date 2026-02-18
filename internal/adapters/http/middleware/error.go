package middleware

import (
	"errors"
	"net/http"

	common "github.com/evythrossell/account-management-api/internal/core/common"
	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var de *common.DomainError

			if errors.As(err, &de) {
				c.JSON(de.HTTPStatusCode(), gin.H{
					"code":    de.Code,
					"message": de.PublicMessage(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "an unexpected error occurred",
				})
			}

			c.Abort()
		}
	}
}
