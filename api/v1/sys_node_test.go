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

func TestBatchDeleteNodeByIds(t *testing.T) {
	r := tests2.GetRouter()
	r.DELETE("/node/delete/batch", BatchDeleteNodeByIds)

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
			url:  "/node/delete/batch?ids=a",
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name: "success",
			s:    &s,
			url:  "/node/delete/batch?ids=1",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "address"}).AddRow(1, "12.34.56.78"))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_node`").WithArgs(sqlmock.AnyArg(), tests2.AnyTime{}, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				// soft delete is update sql
				mock.ExpectExec("UPDATE `tb_sys_node`").WithArgs(tests2.AnyTime{}, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
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

func TestCreateNode(t *testing.T) {
	r := tests2.GetRouter()
	r.POST("/node/create", CreateNode)

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
			url:   "/node/create",
			param: `{"address": "10.34.23.57", "sshPort":"aa"}`,
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name:  "success",
			s:     &s,
			url:   "/node/create",
			param: `{"address": "10.34.23.57", "sshPort":22, "labelIds": [1]}`,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs("10.34.23.57").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_node`").
					WillReturnResult(sqlmock.NewResult(1, 1))
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

func TestGetNodes(t *testing.T) {
	r := tests2.GetRouter()
	r.GET("/node/list", GetNodes)

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
			url:  "/node/list?aa=12.34",
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name: "success",
			s:    &s,
			url:  "/node/list?address=12.34&asset=aa&creator=12345678&health=1&pageNum=1&pageSize=10&noPagination=true",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(
					fmt.Sprintf("%%%s%%", "12.34"),
					fmt.Sprintf("%%%s%%", "aa"),
					fmt.Sprintf("%%%s%%", "12345678"),
					1,
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

func TestRefreshNodeInfo(t *testing.T) {
	r := tests2.GetRouter()
	r.PATCH("/node/refresh/:nodeId", RefreshNodeInfo)

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
			url:  "/node/refresh/a",
			invoke: func() {
			},
			respCode: 405,
		},
		{
			name: "success",
			s:    &s,
			url:  "/node/refresh/1",
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
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
