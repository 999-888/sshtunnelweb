package middleware

import (
	"github.com/gin-gonic/gin"
)

func Tt() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("new-token", "tttoken")
		c.Header("new-expires-at", "20220519")
		c.Next()
	}
}
