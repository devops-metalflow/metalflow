package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const Options = "OPTIONS"

// Cors 处理跨域请求，支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin == "null" || origin == "" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token,"+
			" Content-Length, Authorization, Token, api-idempotence-token, x-tenant-id, x-account-id, x-auth-value,"+
			" x-emp-no, Accept, X-Lang-Id, pinpoint-flags, pinpoint-host, pinpoint-pAppName, pinpoint-pAppType,"+
			" pinpoint-pSpanID, pinpoint-SpanID, pinpoint-TraceID, apm_biz_key")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin,"+
			" Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == Options {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
