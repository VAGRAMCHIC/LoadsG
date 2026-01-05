package router

import (
	"github.com/gin-gonic/gin"

	"loadsg/lib/handler"
)

func registerAuth(r *gin.RouterGroup, 
	h *handler.Handler,
	
){
	auth := r.Group("/auth")
	auth.POST("/login", h.Auth.Login)
	auth.POST("/refresh", h.Auth.Refresh)
}
