package v1

import (
	"encoding/json"
	"errors"
	"fmt"
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

func TestBatchDeleteApiByIds(t *testing.T) {
	r := tests2.GetRouter()
	r.DELETE("/api/delete/batch", BatchDeleteApiByIds)

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
			url:  "/api/delete/batch?ids=a",
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name:  "success",
			s:     &s,
			url:   "/api/delete/batch?ids=1",
			param: "",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_api`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			respCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			tests2.SetCasbinEnforcer("../../tests/rbac_model.conf")
			// mock一个http请求
			req := httptest.NewRequest(
				http.MethodDelete,           // 请求方法
				tt.url,                      // 请求URL
				strings.NewReader(tt.param), // 请求参数
			)
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

func TestCreateApi(t *testing.T) {
	r := tests2.GetRouter()
	r.POST("/api/create", CreateApi)

	mock := tests2.GetMock()
	s := service.New(nil)

	tests2.MockGetCurrentUser(mock, uint(0))

	tests := []struct {
		name     string
		s        *service.MysqlService
		url      string
		param    string
		invoke   func()
		respCode int
	}{
		{
			name:  "fail",
			s:     &s,
			url:   "/api/create",
			param: `{"method": "2"}`,
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name:  "success",
			s:     &s,
			url:   "/api/create",
			param: `{"method":"POST","path":"/v1/test","category":"test"}`,
			invoke: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_api`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			respCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			tests2.SetCasbinEnforcer("../../tests/rbac_model.conf")
			tests2.SetLog()
			tests2.SetValidate()
			// mock一个http请求
			req := httptest.NewRequest(
				http.MethodPost,             // 请求方法
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

func TestGetAllApiGroupByCategoryByRoleId(t *testing.T) {
	r := tests2.GetRouter()
	r.GET("/api/all/category/:roleId", GetAllApiGroupByCategoryByRoleId)

	mock := tests2.GetMock()
	s := service.New(nil)

	tests2.MockGetCurrentUser(mock, uint(0))

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
			url:  "/api/all/category/1",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WillReturnError(errors.New("DB error"))
			},
			respCode: 405,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			tests2.SetCasbinEnforcer("../../tests/rbac_model.conf")
			tests2.SetLog()
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

func TestGetApis(t *testing.T) {
	r := tests2.GetRouter()
	r.GET("/api/list", GetApis)

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
			url:  "/api/list?meth=1",
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name: "success",
			s:    &s,
			url:  "/api/list?method=test&path=v1/test&category=aa&creator=12345678&pageNum=1&pageSize=10&noPagination=true",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(
					fmt.Sprintf("%%%s%%", "test"),
					fmt.Sprintf("%%%s%%", "v1/test"),
					fmt.Sprintf("%%%s%%", "aa"),
					fmt.Sprintf("%%%s%%", "12345678"),
				).WillReturnRows(sqlmock.NewRows([]string{"id"}))
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

func TestUpdateApiById(t *testing.T) {
	r := tests2.GetRouter()
	r.PATCH("/api/update/:apiId", UpdateApiById)

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
			url:  "/api/update/a",
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name:  "success",
			s:     &s,
			url:   "/api/update/1",
			param: `{"method":"POST"}`,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_api`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			respCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			tests2.SetCasbinEnforcer("../../tests/rbac_model.conf")
			tests2.SetLog()
			tests2.SetValidate()
			// mock一个http请求
			req := httptest.NewRequest(
				http.MethodPatch,            // 请求方法
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
