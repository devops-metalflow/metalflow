package service

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"metalflow/pkg/request"
	tests2 "metalflow/tests"
	"testing"
)

func TestMysqlService_GetRoleIdsBySort(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		currentRoleSort uint
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
			args: args{currentRoleSort: uint(1)},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.currentRoleSort).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{currentRoleSort: uint(1)},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.currentRoleSort).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if _, err := tt.s.GetRoleIdsBySort(tt.args.currentRoleSort); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetRoleIdsBySort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_GetRoles(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)

	n := uint(1)
	type args struct {
		req *request.RoleListRequestStruct
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
			args: args{req: &request.RoleListRequestStruct{
				Name:    "tester",
				Keyword: "tester",
				Creator: "jack",
				Status:  &n,
			}},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(
					sqlmock.AnyArg(),
					fmt.Sprintf("%%%s%%", args2.req.Name),
					fmt.Sprintf("%%%s%%", args2.req.Keyword),
					fmt.Sprintf("%%%s%%", args2.req.Creator),
					args2.req.Status,
				).WillReturnError(errors.New("sql error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{req: &request.RoleListRequestStruct{
				Name:    "tester",
				Keyword: "tester",
				Creator: "jack",
				Status:  &n,
			}},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(
					sqlmock.AnyArg(),
					fmt.Sprintf("%%%s%%", args2.req.Name),
					fmt.Sprintf("%%%s%%", args2.req.Keyword),
					fmt.Sprintf("%%%s%%", args2.req.Creator),
					args2.req.Status,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if _, err := tt.s.GetRoles(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetRoles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_UpdateRoleApisById(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		id  uint
		req request.UpdateIncrementalIdsRequestStruct
	}
	tests := []struct {
		name    string
		s       *MysqlService
		args    args
		invoke  func(args2 args)
		wantErr bool
	}{
		{
			name: "fail0",
			s:    &s,
			args: args{
				id: 1,
				req: request.UpdateIncrementalIdsRequestStruct{
					Create: []uint{1},
					Delete: []uint{2},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.id).WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "fail2",
			s:    &s,
			args: args{
				id: 1,
				req: request.UpdateIncrementalIdsRequestStruct{
					Create: []uint{1},
					Delete: []uint{2},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(args2.req.Delete[0]).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail3",
			s:    &s,
			args: args{
				id: 1,
				req: request.UpdateIncrementalIdsRequestStruct{
					Create: []uint{1},
					Delete: []uint{2},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(args2.req.Delete[0]).WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(args2.req.Create[0]).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				id: 1,
				req: request.UpdateIncrementalIdsRequestStruct{
					Create: []uint{1},
					Delete: []uint{2},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(args2.req.Delete[0]).WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_api`").WithArgs(args2.req.Create[0]).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			tests2.SetCasbinEnforcer("../../tests/rbac_model.conf")
			if err := s.UpdateRoleApisById(tt.args.id, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRoleApisById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
