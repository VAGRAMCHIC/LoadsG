// internal/router/router.go
package router

import (
	"loadsg/lib/handler"
	"loadsg/lib/security"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	h *handler.Handler,
	s *security.JWTManager,
) {
	v1 := r.Group("/v1")
	registerAuth(v1, h)
	registerUserPrivate(v1, h, s)
	registerLoadManagerPrivate(v1, h, s)
}
