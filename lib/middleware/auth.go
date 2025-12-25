package middleware

import "github.com/gin-gonic/gin"

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{
                "error": "authorization required",
            })
            return
        }
        c.Next()
    }
}

