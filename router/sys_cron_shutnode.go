package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "metalflow/api/v1"
)

func InitCronShutNodeRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/cron")
	{ // nolint:gocritic
		router1.GET("/list", v1.GetCronShutTasks)
		router1.POST("/create", v1.CreateCronShutTask)
		// 为了支持小程序，临时更新为POST请求方式
		router1.POST("/update/:shutId", v1.UpdateCronShutTaskById)
		router1.DELETE("/delete/batch", v1.BatchDeleteCronShutTask)
	}
	return r
}
