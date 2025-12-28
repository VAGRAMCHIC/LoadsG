package router


import (
	"github.com/gin-gonic/gin"

	"loadsg/lib/handler"
	"loadsg/lib/middleware"
)

func registerUserPrivate(r *gin.RouterGroup, 
	h *handler.Handler,
	
){
	uR := r.Group("/users")

	uR.Use(middleware.AuthRequired())
	r.POST("/create", h.User.Create)
}
