package service

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	tests2 "metalflow/tests"
	"reflect"
	"testing"
)

// nolint: funlen
func TestMysqlService_Create(t *testing.T) {
	mock := tests2.GetMock()
	// 创建服务
	s := New(nil)

	type args struct {
		req  interface{}
		user interface{}
	}
	tests := []struct {
		name    string
		s       *MysqlService
		args    args
		invoke  func(args)
		wantErr bool
	}{
		{
			name: "create-success",
			s:    &s,
			args: args{
				req: request.CreateUserRequestStruct{
					Username: "001",
					Password: "001",
				},
				user: &models.SysUser{},
			},
			invoke: func(args args) {
				// 修改
				mock.ExpectBegin() // mock事务的开始
				mock.ExpectExec("INSERT INTO `tb_sys_user`").WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
					args.req.(request.CreateUserRequestStruct).Username,
					args.req.(request.CreateUserRequestStruct).Password,
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit() // mock事务的结束
			},
			wantErr: false,
		},
		{
			name: "create-failed",
			s:    &s,
			args: args{
				req: request.CreateUserRequestStruct{
					Username: "002",
					Password: "002",
				},
				user: &models.SysUser{},
			},
			invoke: func(args args) {
				mock.ExpectBegin() // mock事务的开始
				mock.ExpectExec("INSERT INTO `tb_sys_user`").WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
					args.req.(request.CreateUserRequestStruct).Username,
					args.req.(request.CreateUserRequestStruct).Password,
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				).WillReturnError(errors.New("username must be unique"))
				mock.ExpectRollback() // 抛出异常，需要回滚
			},
			wantErr: true,
		},
		{
			name: "non-pointer-model",
			s:    &s,
			args: args{
				req: request.CreateUserRequestStruct{
					Username: "003",
					Password: "003",
				},
				user: models.SysUser{},
			},
			invoke: func(args args) {
				mock.ExpectExec("INSERT INTO `tb_sys_user`").WithArgs(
					args.req.(request.CreateUserRequestStruct).Username,
					args.req.(request.CreateUserRequestStruct).Password,
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "create-slice",
			s:    &s,
			args: args{
				req: []request.CreateUserRequestStruct{
					{Username: "004", Password: "004"},
					{Username: "005", Password: "005"},
				},
				user: &models.SysUser{},
			},
			invoke: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_user`").WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
					args.req.([]request.CreateUserRequestStruct)[0].Username,
					args.req.([]request.CreateUserRequestStruct)[0].Password,
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
					args.req.([]request.CreateUserRequestStruct)[1].Username,
					args.req.([]request.CreateUserRequestStruct)[1].Password,
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(3, 2))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "create-slice-slice",
			s:    &s,
			args: args{
				req: []request.CreateUserRequestStruct{
					{Username: "006", Password: "006"},
					{Username: "007", Password: "007"},
				},
				user: []models.SysUser{},
			},
			invoke: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_user`").WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
					args.req.([]request.CreateUserRequestStruct)[0].Username,
					args.req.([]request.CreateUserRequestStruct)[0].Password,
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
					args.req.([]request.CreateUserRequestStruct)[1].Username,
					args.req.([]request.CreateUserRequestStruct)[1].Password,
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(5, 2))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "non-pointer-model" {
				err := tt.s.Create(tt.args.req, tt.args.user)
				assert.Equal(t, err.Error(), "model must be a pointer")
				return
			}

			tt.invoke(tt.args)
			if err := tt.s.Create(tt.args.req, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_DeleteByIds(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)

	type MyModel struct {
		Id   uint
		Name string
	}
	type args struct {
		ids   []uint
		model any
	}
	tests := []struct {
		name    string
		s       *MysqlService
		args    args
		invoke  func()
		wantErr bool
	}{
		{
			name: "success",
			s:    &s,
			args: args{
				ids:   []uint{1},
				model: &MyModel{},
			},
			invoke: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `my_models`").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "fail",
			s:    &s,
			args: args{
				ids:   []uint{1},
				model: &MyModel{},
			},
			invoke: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `my_models`").WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			if err := tt.s.DeleteByIds(tt.args.ids, tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.DeleteByIds() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// nolint: funlen
func TestMysqlService_Find(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	tests2.SetLog()
	countCache := true

	type args struct {
		query *gorm.DB
		page  *response.PageInfo
		model any
	}
	tests := []struct {
		name    string
		s       *MysqlService
		invoke  func()
		args    args
		wantErr bool
	}{
		{
			name: "non-pointer-model",
			s:    &s,
			args: args{
				model: nil,
			},
			invoke:  func() {},
			wantErr: true,
		},
		{
			name: "no-pagination",
			s:    &s,
			args: args{
				query: s.DB.Model(&models.SysUser{}),
				page: &response.PageInfo{
					CountCache:   &countCache,
					NoPagination: true,
				},
				model: &[]models.SysUser{},
			},
			invoke: func() {
				// 查询
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
		{
			name: "pagination-non-skipcount-set-cache",
			s:    &s,
			args: args{
				query: s.DB.Model(&models.SysUser{}),
				page: &response.PageInfo{
					NoPagination: false,
					SkipCount:    false,
				},
				model: &[]models.SysUser{},
			},
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
			wantErr: false,
		},
		{
			name: "pagination-non-skipcount-get-cache",
			s:    &s,
			args: args{
				query: s.DB.Model(&models.SysUser{}),
				page: &response.PageInfo{
					NoPagination: false,
					SkipCount:    false,
					CountCache:   &countCache,
				},
				model: &[]models.SysUser{},
			},
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
			wantErr: false,
		},
		{
			name: "pagination-skipcount-non-limitprimary",
			s:    &s,
			args: args{
				query: s.DB.Model(&models.SysUser{}),
				page: &response.PageInfo{
					NoPagination: false,
					SkipCount:    true,
					Total:        int64(1),
				},
				model: &[]models.SysUser{},
			},
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
			wantErr: false,
		},
		{
			name: "pagination-skipcount-limitprimary",
			s:    &s,
			args: args{
				query: s.DB.Model(&models.SysUser{}),
				page: &response.PageInfo{
					NoPagination: false,
					SkipCount:    true,
					Total:        int64(1),
					LimitPrimary: "id",
				},
				model: &[]models.SysUser{},
			},
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
			wantErr: false,
		},
		{
			name: "pagination-skipcount-limitprimary-invalidmodel",
			s:    &s,
			args: args{
				query: s.DB.Model(""),
				page: &response.PageInfo{
					NoPagination: false,
					SkipCount:    true,
					Total:        int64(1),
					LimitPrimary: "id",
				},
				model: &[]models.SysUser{},
			},
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `sys_users`").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tests2.SetCasbinEnforcer("../../tests/rbac_model.conf")
			tt.invoke()
			if err := tt.s.Find(tt.args.query, tt.args.page, tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_UpdateById(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type Req struct {
		Name string
	}
	type MyModel struct {
		Id   uint
		Name string
	}

	type args struct {
		id    uint
		req   any
		model any
	}
	tests := []struct {
		name    string
		s       *MysqlService
		args    args
		invoke  func()
		wantErr bool
	}{
		{
			name: "non_pointer1",
			s:    &s,
			args: args{
				model: nil,
			},
			invoke:  func() {},
			wantErr: true,
		},
		{
			name: "record_not_found",
			s:    &s,
			args: args{
				id:    uint(1),
				req:   &Req{"test"},
				model: &MyModel{},
			},
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `my_models`").WithArgs(1).WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "ok",
			s:    &s,
			args: args{
				id:    uint(1),
				req:   &Req{"test"},
				model: &MyModel{},
			},
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `my_models`").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `my_models`").WithArgs("test", 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			if err := tt.s.UpdateById(tt.args.id, tt.args.req, tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("UpdateById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want MysqlService
	}{
		{
			name: "base",
			args: args{nil},
			want: MysqlService{TX: global.Mysql, DB: global.Mysql},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
