package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "metalflow/api/v1"
)

// InitApiRouter 接口路由
func InitApiRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/api")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/api")
	{ // nolint:gocritic
		router1.GET("/list", v1.GetApis)
		router1.GET("/all/category/:roleId", v1.GetAllApiGroupByCategoryByRoleId)
		router2.POST("/create", v1.CreateApi)
		router1.PATCH("/update/:apiId", v1.UpdateApiById)
		router1.DELETE("/delete/batch", v1.BatchDeleteApiByIds)
	}
	return r
}
