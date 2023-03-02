package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if token := extractJWT(auth); token == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func extractJWT(auth string) *string {
	if strings.HasPrefix(auth, "Bearer ") {
		token := auth[7:]
		return &token
	}

	return nil
}
