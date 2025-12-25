package middleware

import (
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

func LoginRateLimit(max int, window time.Duration) gin.HandlerFunc {
    var mu sync.Mutex
    attempts := make(map[string]int)

    return func(c *gin.Context) {
        ip := c.ClientIP()

        mu.Lock()
        attempts[ip]++
        count := attempts[ip]
        mu.Unlock()

        if count > max {
            c.AbortWithStatusJSON(429, gin.H{
                "error": "too many login attempts",
            })
            return
        }

        go func() {
            time.Sleep(window)
            mu.Lock()
            attempts[ip]--
            mu.Unlock()
        }()

        c.Next()
    }
}
