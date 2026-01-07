package router

import (
	"loadsg/lib/handler"
	"loadsg/lib/middleware"
	"loadsg/lib/security"

	"github.com/gin-gonic/gin"
)

func registerUserPrivate(
	r *gin.RouterGroup,
	h *handler.Handler,
	s *security.JWTManager,
) {
	uR := r.Group("/users")
	uR.Use(middleware.AuthRequired(s))

	uR.POST("/create", h.User.Create)
}
