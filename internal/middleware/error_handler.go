package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/shennawardana23/graphql-pba/internal/util/exception"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			if validationErr, ok := err.(*exception.CustomError); ok {
				c.JSON(400, gin.H{
					"errors": []gin.H{
						{
							"message": validationErr.Message,
							"extensions": gin.H{
								"code":    validationErr.Code,
								"details": validationErr.Details,
							},
						},
					},
				})
				return
			}
			// Handle other types of errors...
		}
	}
}
