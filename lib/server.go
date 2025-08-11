package lib

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
)

// var jwtKey = []byte("oquooKiezee6ohy") // Лучше хранить в .env
var users = map[string]string{} // username -> password (для примера, без БД)

// Claims описывает содержимое токена
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Server(jwtKey []byte) {
	r := gin.Default()

	// Регистрация нового пользователя
	r.POST("/register", func(c *gin.Context) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		if _, exists := users[creds.Username]; exists {
			c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
			return
		}
		users[creds.Username] = creds.Password
		c.JSON(http.StatusOK, gin.H{"message": "registered"})
	})

	// Авторизация (выдача токена)
	r.POST("/login", func(c *gin.Context) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		if users[creds.Username] != creds.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// Генерируем токен
		expirationTime := time.Now().Add(1 * time.Hour)
		claims := &Claims{
			Username: creds.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})

	// Группа защищённых роутов
	auth := r.Group("/api")
	auth.Use(JWTAuthMiddleware(jwtKey))
	{
		auth.GET("/me", func(c *gin.Context) {
			user, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{"user": user})
		})
		auth.POST("/load", func(c *gin.Context) {
			var loadRequest LoadRrequest

			headers := map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
				"User-Agent":   "LoadsG/1.0",
			}
			var httpLoadRequest HTTPLoadRequest
			httpLoadRequest.Id = 0
			httpLoadRequest.HttpHead = CreateHttpHead("GET", "http://test.customlabs.ru/test2/", "HTTP/1.1", headers)
			httpLoadRequest.Body = "test"

			if err := c.ShouldBindBodyWithJSON(&loadRequest); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
				return
			}
			log.Print(loadRequest)
			request, host := BuildHttpRequest(httpLoadRequest.HttpHead, httpLoadRequest.Body)

			requests := make([]string, loadRequest.Amount)
			for i := range loadRequest.Amount {
				requests[i] = request
			}
			var wg sync.WaitGroup
			for _, req := range requests {
				wg.Add(1)
				go func(r string) {
					defer wg.Done()
					SendHttpRequest(r, host)
					log.Print(r, host)
				}(req)
			}
			user, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"user":          user,
				"status":        "ok",
				"sent_requests": loadRequest.Amount,
			})

		})
	}

	r.Run(":8080")
}

/*
	func JWTAuthMiddleware() gin.HandlerFunc {
		return func(c *gin.Context) {
			tokenString := c.GetHeader("Authorization")
			if tokenString == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			fmt.Print(token, err)
			if err != nil || !token.Valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}

			c.Set("username", claims.Username)
			c.Next()
		}
	}
*/
func JWTAuthMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			tokenString = strings.TrimSpace(tokenString)
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
