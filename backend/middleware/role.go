package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("user_role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authenticated"})
			return
		}
		role, _ := v.(string)
		if role != required {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
