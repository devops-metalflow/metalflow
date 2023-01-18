package router

import (
	v1 "metalflow/api/v1"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitUserRouter 用户路由
func InitUserRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/user")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/user")
	{ // nolint:gocritic
		router1.GET("/info", v1.GetUserInfo)
		router1.GET("/list", v1.GetUsers)
		router2.POST("/create", v1.CreateUser)
		router1.PATCH("/update/:userId", v1.UpdateUserById)
		router1.DELETE("/delete/batch", v1.BatchDeleteUserByIds)
	}
	return r
}
