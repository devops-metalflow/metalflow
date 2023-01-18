package router

import (
	v1 "metalflow/api/v1"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitMenuRouter 菜单路由
func InitMenuRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/menu")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/menu")
	{ // nolint:gocritic
		router1.GET("/tree", v1.GetMenuTree)
		router1.GET("/all/:roleId", v1.GetAllMenuByRoleId)
		router1.GET("/list", v1.GetMenus)
		router2.POST("/create", v1.CreateMenu)
		router1.PATCH("/update/:menuId", v1.UpdateMenuById)
		router2.DELETE("/delete/batch", v1.BatchDeleteMenuByIds)
	}
	return r
}
