package lib

import (
	"net/http"
	"strings"
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
		auth.POST("/load/http", func(c *gin.Context) {
			var loadRequest HTTPLoadRequest
			if err := c.BindJSON(&loadRequest); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			result, _ := RunLoad(loadRequest)
			c.JSON(http.StatusOK, gin.H{
				"status":        "ok",
				"sent_requests": result,
			})
		})
	}

	r.Run(":8080")
}

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
