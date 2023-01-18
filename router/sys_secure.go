package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "metalflow/api/v1"
)

// InitSecureRouter 系统守护路由
func InitSecureRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router := GetCasbinRouter(r, authMiddleware, "/secure")
	{ // nolint:gocritic
		router.GET("risk/:nodeId", v1.GetRiskCountById)
		router.GET("/stats/:nodeId", v1.GetNodeImagesById)
		router.POST("/container-report/:nodeId", v1.GetNodeImageSecureById)
		router.POST("/bare-report/:nodeId", v1.GetNodeBareSecureById)
		router.GET("/score/:nodeId", v1.GetNodeSecurityScoreById)
		router.POST("/bare/:nodeId", v1.RunNodeBareSecurityById)
		router.POST("/container/:nodeId", v1.RunNodeContainerSecurityById)
		router.POST("/fix/:nodeId", v1.FixNodeSecurityIssuesById)
	}
	return r
}
