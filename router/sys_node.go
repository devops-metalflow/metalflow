package router

import (
	v1 "metalflow/api/v1"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitMachineRouter 机器节点路由
func InitMachineRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (i gin.IRoutes) {
	router1 := GetCasbinRouter(r, authMiddleware, "/node")
	router2 := GetCasbinAndIdempotenceRouter(r, authMiddleware, "/node")
	{ // nolint:gocritic
		router1.POST("/shell/connect", v1.NodeConnect)
		router1.GET("/shell/ws", v1.NodeShellWs)
		router1.GET("/vnc/ws", v1.NodeVncWs)
		router1.PATCH("/shell/ws/resize", v1.ResizeWs)
		router1.GET("/shell/dir", v1.GetSshDirInfo)
		router1.GET("/shell/file", v1.GetSshFile)
		router1.POST("/shell/file/download", v1.DownloadFile)
		router1.PATCH("/shell/file/update", v1.UpdateFile)
		router1.GET("/list", v1.GetNodes)
		router2.POST("/create", v1.CreateNode)
		router1.POST("/reboot", v1.BatchRebootNodeByIds)
		router1.PATCH("/update/:nodeId", v1.UpdateNodeById)
		router1.DELETE("/delete/batch", v1.BatchDeleteNodeByIds)
		router1.PATCH("/refresh/:nodeId", v1.RefreshNodeInfo)
	}
	return r
}
