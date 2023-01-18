package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "metalflow/api/v1"
)

// InitOperationLogRouter for handle operation log
func InitOperationLogRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router := GetCasbinRouter(r, authMiddleware, "/operation/log")
	{ // nolint:gocritic
		router.GET("/list", v1.GetOperationLogs)
		router.DELETE("/delete/batch", v1.BatchDeleteOperationLogByIds)
	}
	return r
}
