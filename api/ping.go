package api

import (
	"metalflow/pkg/response"

	"github.com/gin-gonic/gin"
)

// Ping 检查服务器是否通畅
func Ping(c *gin.Context) {
	response.SuccessWithData("pong")
}
