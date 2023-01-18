package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "metalflow/api/v1"
)

// InitDashboardRouter 接口路由
func InitDashboardRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/dashboard")
	{ // nolint:gocritic
		router1.GET("/countData", v1.GetCountData)
		router1.GET("/regionNodeData", v1.GetRegionNodeData)
		router1.GET("/managerNodeData", v1.GetManagerNodeData)
		router1.GET("/performanceNodeData", v1.GetPerformanceNodeData)
		router1.GET("/healthNodeData", v1.GetHealthNodeData)
	}
	return r
}
