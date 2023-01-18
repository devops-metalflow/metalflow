package service

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"metalflow/models"
	"metalflow/pkg/request"
	tests2 "metalflow/tests"
	"testing"
)

func TestMysqlService_CreateApi(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		req *request.CreateApiRequestStruct
	}
	tests := []struct {
		name    string
		s       *MysqlService
		args    args
		invoke  func(args2 args)
		wantErr bool
	}{
		{
			name: "fail",
			s:    &s,
			args: args{
				req: &request.CreateApiRequestStruct{
					Method:   "POST",
					Path:     "/v1/test",
					Category: "test",
					Title:    "test",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_api`").WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "fail2",
			s:    &s,
			args: args{
				req: &request.CreateApiRequestStruct{
					Method:   "GET",
					Path:     "/v1/test2",
					Category: "test",
					RoleIds:  []uint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_api`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(1).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				req: &request.CreateApiRequestStruct{
					Method:   "GET",
					Path:     "/v1/test2",
					Category: "test",
					RoleIds:  []uint{2},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_api`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			tests2.SetCasbinEnforcer("../../tests/rbac_model.conf")
			if err := tt.s.CreateApi(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.CreateApi() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_DeleteApiByIds(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		ids []uint
	}
	tests := []struct {
		name    string
		s       *MysqlService
		args    args
		invoke  func(args2 args)
		wantErr bool
	}{
		{
			name: "failSearch",
			s:    &s,
			args: args{
				ids: []uint{1},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(args2.ids[0]).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "failDelete",
			s:    &s,
			args: args{
				ids: []uint{1},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(args2.ids[0]).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				// soft delete is update sql
				mock.ExpectExec("UPDATE `tb_sys_api`").WithArgs(tests2.AnyTime{}, args2.ids[0]).WillReturnError(errors.New("DB delete error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				ids: []uint{1},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(args2.ids[0]).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				// soft delete is update sql
				mock.ExpectExec("UPDATE `tb_sys_api`").WithArgs(tests2.AnyTime{}, args2.ids[0]).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			tests2.SetCasbinEnforcer("../../tests/rbac_model.conf")
			if err := tt.s.DeleteApiByIds(tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.DeleteApiByIds() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_GetAllApiGroupByCategoryByRoleId(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)

	sort := uint(2)
	type args struct {
		currentRole *models.SysRole
		roleId      uint
	}
	tests := []struct {
		name    string
		s       *MysqlService
		args    args
		invoke  func(args2 args)
		wantErr bool
	}{
		{
			name: "fail",
			s:    &s,
			args: args{
				currentRole: &models.SysRole{
					Model:   models.Model{},
					Name:    "测试员",
					Keyword: "tester",
					Sort:    &sort,
				},
				roleId: 1,
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail2",
			s:    &s,
			args: args{
				currentRole: &models.SysRole{
					Model:   models.Model{},
					Name:    "测试员",
					Keyword: "tester",
					Sort:    &sort,
				},
				roleId: 1,
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.roleId).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				currentRole: &models.SysRole{
					Model:   models.Model{},
					Name:    "测试员",
					Keyword: "tester",
					Sort:    &sort,
				},
				roleId: 1,
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.roleId).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			tests2.SetCasbinEnforcer("../../tests/rbac_model.conf")
			if _, _, err := tt.s.GetAllApiGroupByCategoryByRoleId(tt.args.currentRole, tt.args.roleId); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetAllApiGroupByCategoryByRoleId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_GetApis(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		req *request.ApiRequestStruct
	}
	tests := []struct {
		name    string
		s       *MysqlService
		args    args
		invoke  func(args2 args)
		wantErr bool
	}{
		{
			name: "fail",
			s:    &s,
			args: args{
				req: &request.ApiRequestStruct{
					Method:   "POST",
					Path:     "/v1/test",
					Category: "test",
					Creator:  "tester",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(
					fmt.Sprintf("%%%s%%", args2.req.Method),
					fmt.Sprintf("%%%s%%", args2.req.Path),
					fmt.Sprintf("%%%s%%", args2.req.Category),
					fmt.Sprintf("%%%s%%", args2.req.Creator),
				).WillReturnError(errors.New("search tb_sys_api error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				req: &request.ApiRequestStruct{
					Method:   "POST",
					Path:     "/v1/test",
					Category: "test",
					Creator:  "tester",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(
					fmt.Sprintf("%%%s%%", args2.req.Method),
					fmt.Sprintf("%%%s%%", args2.req.Path),
					fmt.Sprintf("%%%s%%", args2.req.Category),
					fmt.Sprintf("%%%s%%", args2.req.Creator),
				).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if _, err := tt.s.GetApis(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetApis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_UpdateApiById(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)

	strValue := "test"
	type args struct {
		id  uint
		req request.UpdateApiRequestStruct
	}
	tests := []struct {
		name    string
		s       *MysqlService
		args    args
		invoke  func(args2 args)
		wantErr bool
	}{
		{
			name: "fail1",
			s:    &s,
			args: args{
				id: 1,
				req: request.UpdateApiRequestStruct{
					Method:   &strValue,
					Path:     &strValue,
					Category: &strValue,
					Desc:     &strValue,
					Title:    &strValue,
					Creator:  &strValue,
					RoleIds:  []uint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api").WithArgs(args2.id).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail2",
			s:    &s,
			args: args{
				id: 1,
				req: request.UpdateApiRequestStruct{
					Method:   &strValue,
					Path:     &strValue,
					Category: &strValue,
					Desc:     &strValue,
					Title:    &strValue,
					Creator:  &strValue,
					RoleIds:  []uint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api").WithArgs(args2.id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_api`").WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				id: 1,
				req: request.UpdateApiRequestStruct{
					Method:   &strValue,
					Path:     &strValue,
					Category: &strValue,
					Desc:     &strValue,
					Title:    &strValue,
					Creator:  &strValue,
					RoleIds:  []uint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api").WithArgs(args2.id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_api`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			tests2.SetCasbinEnforcer("../../tests/rbac_model.conf")
			if err := tt.s.UpdateApiById(tt.args.id, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.UpdateApiById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
