package router

import (
	"github.com/gin-gonic/gin"

	"loadsg/lib/handler"
)

func registerAuth(r *gin.RouterGroup, 
	h *handler.Handler,
	
){
	r.POST("/login", h.Auth.Login)
}
