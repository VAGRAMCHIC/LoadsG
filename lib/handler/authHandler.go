package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"loadsg/lib/dto"
	"loadsg/lib/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(400)
		return
	}

	pair, err := h.service.Login(
		c.Request.Context(),
		req.UID,
		req.TokenHash,
	)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	setRefreshCookie(c, pair.RefreshToken, pair.RefreshExp)

	c.JSON(200, gin.H{
		"access_token": pair.AccessToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	pair, err := h.service.Refresh(c.Request.Context(), refresh)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	setRefreshCookie(c, pair.RefreshToken, pair.RefreshExp)

	c.JSON(200, gin.H{
		"access_token": pair.AccessToken,
	})
}

func setRefreshCookie(c *gin.Context, token string, exp time.Time) {
	c.SetCookie(
		"refresh_token",
		token,
		int(time.Until(exp).Seconds()),
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)
}
