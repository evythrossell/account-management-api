package middleware

import (
	"errors"
	"log"
	"net/http"

	common "github.com/evythrossell/account-management-api/pkg"
	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			log.Printf("[DEBUG] Error received: %T - %v", err, err)

			var de *common.DomainError
			if errors.As(err, &de) {
				log.Printf("[DEBUG] Error is DomainError - Code: %s, HTTPStatus: %d", de.Code, de.HTTPStatusCode())
				c.AbortWithStatusJSON(de.HTTPStatusCode(), gin.H{
					"code":    de.Code,
					"message": de.PublicMessage(),
				})
				return
			}
			if errors.Is(err, common.ErrAccountNotFound) {
				log.Printf("[DEBUG] Error is ErrAccountNotFound")
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"code":    "NOT_FOUND_ERROR",
					"message": "account not found",
				})
				return
			}

			log.Printf("[DEBUG] Error is unmapped, returning 500")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "an unexpected error occurred",
			})
		}
	}
}
