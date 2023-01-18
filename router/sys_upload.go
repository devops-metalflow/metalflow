package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "metalflow/api/v1"
)

// InitUploadRouter 文件上传路由
func InitUploadRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router := GetCasbinRouter(r, authMiddleware, "/upload")
	{ // nolint:gocritic
		router.GET("/file", v1.UploadFileChunkExists)
		router.POST("/file", v1.UploadFile)
		router.POST("/merge", v1.UploadMerge)
	}
	return r
}
