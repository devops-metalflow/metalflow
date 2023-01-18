package v1

import (
	"encoding/json"
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

func TestBatchDeleteLabelByIds(t *testing.T) {
	r := tests2.GetRouter()
	r.DELETE("/label/delete/batch", BatchDeleteLabelByIds)

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
			url:  "/label/delete/batch?ids=a",
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name:  "success",
			s:     &s,
			url:   "/label/delete/batch?ids=1",
			param: "",
			invoke: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE (.*)").WithArgs(tests2.AnyTime{}, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			respCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
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

func TestCreateLabel(t *testing.T) {
	r := tests2.GetRouter()
	r.POST("/label/create", CreateLabel)

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
			name:  "fail",
			s:     &s,
			url:   "/label/create",
			param: `{"name": "2"}`,
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name:  "success",
			s:     &s,
			url:   "/label/create",
			param: `{"name":"test"}`,
			invoke: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_label`").WillReturnResult(sqlmock.NewResult(1, 1))
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

func TestGetLabels(t *testing.T) {
	r := tests2.GetRouter()
	r.GET("/label/list", GetLabels)

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
			url:  "/label/list?lab=test",
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name: "success",
			s:    &s,
			url:  "/label/list?name=test&creator=12345678&pageNum=1&pageSize=10&noPagination=true",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(
					fmt.Sprintf("%%%s%%", "test"),
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

func TestUpdateLabelById(t *testing.T) {
	r := tests2.GetRouter()
	r.PATCH("/label/update/:labelId", UpdateLabelById)

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
			url:  "/label/update/a",
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name:  "success",
			s:     &s,
			url:   "/label/update/1",
			param: `{"name":"test"}`,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_label`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			respCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			tests2.SetLog()
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
