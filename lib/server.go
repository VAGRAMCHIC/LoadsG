package lib

import (
	"net/http"
	"strings"
	"time"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/golang-jwt/jwt/v5"
)


// Claims описывает содержимое токена
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Server(jwtKey []byte, maxConcurrent int, pgConn *pgx.Conn) {
	r := gin.Default()

	// Регистрация нового пользователя
	r.POST("/register", func(c *gin.Context) {
		var creds User

		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fcreds, err := GetUser(pgConn, creds.Username)
		if fcreds.Username != "" {
			c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := InsertUser(pgConn, creds); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "registered"})
	})

	// Авторизация (выдача токена)
	r.POST("/login", func(c *gin.Context) {
		var creds User
		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		fcreds, err := GetUser(pgConn, creds.Username)
		if fcreds.Username == "" {
			c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
			return
		}
		fmt.Printf("%s, %s, %d", fcreds.Username, fcreds.Password, fcreds.Id)
		if fcreds.Password != creds.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials", fcreds.Password: creds.Password})
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
			reportURL := "http://test.customlabs.ru/"
			statusCodes, loadDuration, rps, _ := RunLoad(loadRequest, maxConcurrent, reportURL)
			c.JSON(http.StatusOK, gin.H{
				"loadDuration": loadDuration,
				"statusCodes":  statusCodes,
				"rps":          rps,
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
