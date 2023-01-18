package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "metalflow/api/v1"
)

// InitLabelRouter 标签路由
func InitLabelRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/label")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/label")
	{ // nolint:gocritic
		router1.GET("/list", v1.GetLabels)
		router2.POST("/create", v1.CreateLabel)
		router1.PATCH("/update/:labelId", v1.UpdateLabelById)
		router1.DELETE("/delete/batch", v1.BatchDeleteLabelByIds)
	}
	return r
}
