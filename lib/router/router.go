// internal/router/router.go
package router

import (
	"github.com/gin-gonic/gin"

	"loadsg/lib/handler"
)

func RegisterRoutes(
	r *gin.Engine,
	h *handler.Handler,
) {
	v1 := r.Group("/v1")
	registerAuth(v1, h)
	registerUserPrivate(v1, h)
}



