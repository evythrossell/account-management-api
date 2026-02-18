package middleware

import (
	"errors"
	"net/http"

	"github.com/evythrossell/account-management-api/internal/core/domain"
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
					"code":    domain.ErrCodeNotFound,
					"message": domain.ErrMsgAccountNotFound,
				})
				return
			}
			if errors.Is(err, common.ErrTransactionNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"code":    domain.ErrCodeNotFound,
					"message": domain.ErrMsgTransactionNotFound,
				})
				return
			}
			if errors.Is(err, common.ErrInvalidAmount) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"code":    domain.ErrCodeValidation,
					"message": domain.ErrMsgAmountInvalid,
				})
				return
			}
			if errors.Is(err, common.ErrInvalidOperation) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"code":    domain.ErrCodeValidation,
					"message": domain.ErrMsgOperationTypeInvalid,
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    domain.ErrCodeInternalError,
				"message": domain.ErrMsgUnexpectedError,
			})
		}
	}
}
