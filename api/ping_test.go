package api

import (
	"github.com/stretchr/testify/assert"
	tests2 "metalflow/tests" //nolint:depguard
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	r := tests2.GetRouter()
	r.GET("/ping", Ping)

	tests := []struct {
		name string
		code int
	}{
		{
			name: "success",
			code: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock一个http请求
			req := httptest.NewRequest(
				"GET",       // 请求方法
				"/ping",     // 请求URL
				http.NoBody, // 请求参数
			)
			// mock一个响应记录器
			w := httptest.NewRecorder()
			// 让server端处理mock请求并记录返回的响应内容
			r.ServeHTTP(w, req)

			// 校验状态码是否符合预期
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
