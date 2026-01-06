package handler

import (
	"net/http"

	"loadsg/lib/dto"
	"loadsg/lib/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
    service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
    return &UserHandler{service: s}
}

func (h *UserHandler) Create(c *gin.Context) {
    var req dto.CreateUserRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    user, err := h.service.Create(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "failed to create user",
        })
        return
    }

    c.JSON(http.StatusCreated, dto.CreateUserResponse{
        UID:    user.UID,
    })
}

