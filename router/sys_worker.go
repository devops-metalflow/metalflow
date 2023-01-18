package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "metalflow/api/v1"
)

// InitWorkerRouter collects all worker handlers associated with request path.
func InitWorkerRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/worker")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/worker")
	{ // nolint:gocritic
		router1.GET("/list", v1.GetWorkers)
		router2.POST("/create", v1.CreateWorker)
		router1.PATCH("/update/:workerId", v1.UpdateWorkerById)
		router1.DELETE("/delete/batch", v1.BatchDeleteWorkerByIds)
	}
	return r
}
