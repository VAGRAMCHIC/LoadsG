// internal/handler/auth_handler.go
package handler

import (
	"log"
	"net/http"

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
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid request body",
        })
        return
    }

    token, refreshToken, err := h.service.Login(
        c.Request.Context(),
        req.UID,
				req.TokenHash,
    )
    if err != nil {
				log.Printf("\nerror: %s", err.Error())
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "invalid id or token",
        })
        return
    }

    c.JSON(http.StatusOK, dto.LoginResponse{
        AccessToken: token, RefreshToken: refreshToken,
    })
}

