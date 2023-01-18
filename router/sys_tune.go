package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "metalflow/api/v1"
)

// InitTuneRouter  系统调优
func InitTuneRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router := GetCasbinRouter(r, authMiddleware, "/tune")
	{ // nolint:gocritic
		router.GET("/score/:nodeId", v1.GetTuneScoreByNodeId)
		router.GET("/auto/list/:nodeId", v1.GetTuneLogsByNodeId)
		router.POST("/cleanup/:nodeId", v1.Cleanup)
		router.POST("/auto/rollback/:nodeId", v1.Rollback)
		router.POST("/auto/set/:nodeId", v1.Set)
		router.POST("/scene/:nodeId", v1.Scene)
		router.POST("/turbo/:nodeId", v1.Turbo)
		router.DELETE("/auto/delete", v1.BatchDeleteTuneLogByIds)
	}
	return r
}
