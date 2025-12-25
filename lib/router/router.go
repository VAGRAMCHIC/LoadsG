// internal/router/router.go
package router

import (
	"github.com/gin-gonic/gin"

	"loadsg/lib/handler"
)

func RegisterRoutes(
	r *gin.Engine,
	user *handler.UserHandler,
	auth *handler.AuthHandler,
) {
	v1 := r.Group("/v1")
	{
		v1.POST("/register", user.Create)
		v1.POST("/login", auth.Login)
	}
}



