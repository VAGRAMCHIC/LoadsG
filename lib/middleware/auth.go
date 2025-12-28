package middleware

import (
	"strings"
	"loadsg/lib/security"
	"github.com/gin-gonic/gin"
)


// internal/middleware/auth.go
func AuthRequired(jwtManager *security.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")

		claims, err := jwtManager.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// передаём userID дальше
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}


