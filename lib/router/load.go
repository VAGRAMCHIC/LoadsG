package router

import (
	"loadsg/lib/handler"
	"loadsg/lib/middleware"
	"loadsg/lib/security"

	"github.com/gin-gonic/gin"
)

func registerLoadManagerPrivate(
	r *gin.RouterGroup,
	h *handler.Handler,
	s *security.JWTManager,
){
	lM:= r.Group("/manager")
	
	lM.Use(middleware.AuthRequired(s))
	
	lF := lM.Group("/http")

	lF.POST("/create-fixed", h.Load.CreateFixedHttp)
}
