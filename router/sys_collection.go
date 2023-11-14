package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "metalflow/api/v1"
)

func InitCollectionRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/collect")
	{ // nolint:gocritic
		router1.GET("/my", v1.GetMyCollections)
		router1.POST("/my", v1.CreateCollection)
		router1.DELETE("/delete/my", v1.DeleteCollectByNodeId)
		router1.PATCH("/update/:collectionId", v1.UpdateDescriptionById)
	}
	return r
}
