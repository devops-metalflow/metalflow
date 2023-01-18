package router

import (
	"github.com/gin-gonic/gin"
)

// InitPublicRouter 公共路由, 任何人可访问
func InitPublicRouter(r *gin.RouterGroup) (i gin.IRoutes) {
	r.Group("/public")
	{ // nolint:gocritic
	}
	return r
}
