package service

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"metalflow/models"
	"metalflow/pkg/request"
	tests2 "metalflow/tests"
	"testing"
)

func TestMysqlService_AuthCheck(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	tests2.SetConfig("../../tests/config.yml")
	type args struct {
		req *request.UserAuthRequestStruct
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
				req: &request.UserAuthRequestStruct{
					Username: "12345678",
					Password: "12345678",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail2",
			s:    &s,
			args: args{
				req: &request.UserAuthRequestStruct{
					Username: "12345678",
					Password: "12345678",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "fail3",
			s:    &s,
			args: args{
				req: &request.UserAuthRequestStruct{
					Username: "12345678",
					Password: "12345678",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WithArgs(args2.req.Username).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if _, err := tt.s.LoginCheck(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.LoginCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_CreateUser(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		req *request.CreateUserRequestStruct
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
			args: args{req: &request.CreateUserRequestStruct{
				Username: "12345678",
				Password: "jack",
			}},
			invoke: func(args2 args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_user`").WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{req: &request.CreateUserRequestStruct{
				Username: "12345678",
				Password: "jack",
			}},
			invoke: func(args2 args) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_user`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if err := tt.s.CreateUser(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_GetUserById(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		id uint
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
				id: 1,
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WithArgs(args2.id, 1).
					WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				id: 1,
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WithArgs(args2.id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if _, err := tt.s.GetUserById(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetUserById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_GetUsers(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)

	constNum := uint(1)
	type args struct {
		req *request.UserListRequestStruct
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
				req: &request.UserListRequestStruct{
					CurrentRole:  models.SysRole{Sort: &constNum},
					Username:     "12345678",
					OfficeName:   "",
					OrgName:      "",
					Introduction: "",
					Creator:      "",
					Status:       &constNum,
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.req.CurrentRole.Sort).
					WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail2",
			s:    &s,
			args: args{
				req: &request.UserListRequestStruct{
					CurrentRole:  models.SysRole{Sort: &constNum},
					Username:     "12345678",
					OfficeName:   "",
					OrgName:      "",
					Introduction: "",
					Creator:      "",
					Status:       &constNum,
					RoleId:       constNum,
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.req.CurrentRole.Sort).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WithArgs(
					args2.req.RoleId,
					fmt.Sprintf("%%%s%%", args2.req.Username),
					args2.req.Status,
				).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				req: &request.UserListRequestStruct{
					CurrentRole:  models.SysRole{Sort: &constNum},
					Username:     "12345678",
					OfficeName:   "",
					OrgName:      "",
					Introduction: "",
					Creator:      "",
					Status:       &constNum,
					RoleId:       constNum,
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_role`").WithArgs(args2.req.CurrentRole.Sort).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WithArgs(
					args2.req.RoleId,
					fmt.Sprintf("%%%s%%", args2.req.Username),
					args2.req.Status,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if _, err := tt.s.GetUsers(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
