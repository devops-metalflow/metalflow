package middleware

import (
	"fmt"
	"metalflow/pkg/global"
	"metalflow/pkg/response"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Exception 全局异常处理中间件
func Exception(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 将异常写入日志
			global.Log.Error(fmt.Sprintf("[Exception]未知异常: %v\n堆栈信息: %v", err, string(debug.Stack())))
			// 服务器异常
			resp := response.Resp{
				Code:   response.InternalServerError,
				Result: map[string]any{},
				Msg:    response.CustomError[response.InternalServerError],
			}
			// 以json方式写入响应
			response.JSON(c, http.StatusOK, resp)
		}
	}()
	c.Next()
}
