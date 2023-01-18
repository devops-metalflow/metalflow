package v1

import (
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	tests2 "metalflow/tests"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetCountData(t *testing.T) {
	r := tests2.GetRouter()
	r.GET("/dashboard/countData", GetCountData)

	mock := tests2.GetMock()
	s := service.New(nil)

	tests := []struct {
		name     string
		s        *service.MysqlService
		url      string
		param    string
		invoke   func()
		respCode int
	}{
		{
			name: "fail",
			s:    &s,
			url:  "/dashboard/countData",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnError(errors.New("DB error"))
			},
			respCode: 405,
		},
		{
			name: "success",
			s:    &s,
			url:  "/dashboard/countData",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			respCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			// mock一个http请求
			req := httptest.NewRequest(
				http.MethodGet,              // 请求方法
				tt.url,                      // 请求URL
				strings.NewReader(tt.param), // 请求参数
			)
			req.Header.Set("Content-Type", "application/json")
			// mock一个响应记录器
			w := httptest.NewRecorder()
			// 让server端处理mock请求并记录返回的响应内容
			r.ServeHTTP(w, req)

			var resp response.Resp
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Nil(t, err)
			// 校验自定义响应码是否符合预期
			assert.Equal(t, tt.respCode, resp.Code)
		})
	}
}

func TestGetHealthNodeData(t *testing.T) {
	r := tests2.GetRouter()
	r.GET("/dashboard/healthNodeData", GetHealthNodeData)

	mock := tests2.GetMock()
	s := service.New(nil)

	tests := []struct {
		name     string
		s        *service.MysqlService
		url      string
		param    string
		invoke   func()
		respCode int
	}{
		{
			name: "fail",
			s:    &s,
			url:  "/dashboard/healthNodeData",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnError(errors.New("DB error"))
			},
			respCode: 405,
		},
		{
			name: "success",
			s:    &s,
			url:  "/dashboard/healthNodeData",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			respCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			// mock一个http请求
			req := httptest.NewRequest(
				http.MethodGet,              // 请求方法
				tt.url,                      // 请求URL
				strings.NewReader(tt.param), // 请求参数
			)
			req.Header.Set("Content-Type", "application/json")
			// mock一个响应记录器
			w := httptest.NewRecorder()
			// 让server端处理mock请求并记录返回的响应内容
			r.ServeHTTP(w, req)

			var resp response.Resp
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Nil(t, err)
			// 校验自定义响应码是否符合预期
			assert.Equal(t, tt.respCode, resp.Code)
		})
	}
}

func TestGetManagerNodeData(t *testing.T) {
	r := tests2.GetRouter()
	r.GET("/dashboard/managerNodeData", GetManagerNodeData)

	mock := tests2.GetMock()
	s := service.New(nil)

	tests := []struct {
		name     string
		s        *service.MysqlService
		url      string
		param    string
		invoke   func()
		respCode int
	}{
		{
			name: "fail",
			s:    &s,
			url:  "/dashboard/managerNodeData",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnError(errors.New("DB error"))
			},
			respCode: 405,
		},
		{
			name: "success",
			s:    &s,
			url:  "/dashboard/managerNodeData",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			respCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			// mock一个http请求
			req := httptest.NewRequest(
				http.MethodGet,              // 请求方法
				tt.url,                      // 请求URL
				strings.NewReader(tt.param), // 请求参数
			)
			req.Header.Set("Content-Type", "application/json")
			// mock一个响应记录器
			w := httptest.NewRecorder()
			// 让server端处理mock请求并记录返回的响应内容
			r.ServeHTTP(w, req)

			var resp response.Resp
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Nil(t, err)
			// 校验自定义响应码是否符合预期
			assert.Equal(t, tt.respCode, resp.Code)
		})
	}
}

func TestGetPerformanceNodeData(t *testing.T) {
	r := tests2.GetRouter()
	r.GET("/dashboard/performanceNodeData", GetPerformanceNodeData)

	mock := tests2.GetMock()
	s := service.New(nil)

	tests := []struct {
		name     string
		s        *service.MysqlService
		url      string
		param    string
		invoke   func()
		respCode int
	}{
		{
			name: "fail",
			s:    &s,
			url:  "/dashboard/performanceNodeData",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnError(errors.New("DB error"))
			},
			respCode: 405,
		},
		{
			name: "success",
			s:    &s,
			url:  "/dashboard/performanceNodeData",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			respCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			// mock一个http请求
			req := httptest.NewRequest(
				http.MethodGet,              // 请求方法
				tt.url,                      // 请求URL
				strings.NewReader(tt.param), // 请求参数
			)
			req.Header.Set("Content-Type", "application/json")
			// mock一个响应记录器
			w := httptest.NewRecorder()
			// 让server端处理mock请求并记录返回的响应内容
			r.ServeHTTP(w, req)

			var resp response.Resp
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Nil(t, err)
			// 校验自定义响应码是否符合预期
			assert.Equal(t, tt.respCode, resp.Code)
		})
	}
}

func TestGetRegionNodeData(t *testing.T) {
	r := tests2.GetRouter()
	r.GET("/dashboard/regionNodeData", GetRegionNodeData)

	mock := tests2.GetMock()
	s := service.New(nil)

	tests := []struct {
		name     string
		s        *service.MysqlService
		url      string
		param    string
		invoke   func()
		respCode int
	}{
		{
			name: "fail",
			s:    &s,
			url:  "/dashboard/regionNodeData",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnError(errors.New("DB error"))
			},
			respCode: 405,
		},
		{
			name: "success",
			s:    &s,
			url:  "/dashboard/regionNodeData",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			respCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			// mock一个http请求
			req := httptest.NewRequest(
				http.MethodGet,              // 请求方法
				tt.url,                      // 请求URL
				strings.NewReader(tt.param), // 请求参数
			)
			req.Header.Set("Content-Type", "application/json")
			// mock一个响应记录器
			w := httptest.NewRecorder()
			// 让server端处理mock请求并记录返回的响应内容
			r.ServeHTTP(w, req)

			var resp response.Resp
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Nil(t, err)
			// 校验自定义响应码是否符合预期
			assert.Equal(t, tt.respCode, resp.Code)
		})
	}
}
