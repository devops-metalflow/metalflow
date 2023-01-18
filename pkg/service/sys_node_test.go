package service

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	tests2 "metalflow/tests"
	"testing"
)

// nolint: funlen
func TestMysqlService_CreateNode(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		req *request.CreateNodeRequestStruct
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
				req: &request.CreateNodeRequestStruct{
					Address:  "10.23.45.67",
					SshPort:  22,
					LabelIds: []request.ReqUint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.req.Address).
					WillReturnError(errors.New("the machine node already exists, please do not repeat the creation"))
			},
			wantErr: true,
		},
		{
			name: "fail2",
			s:    &s,
			args: args{
				req: &request.CreateNodeRequestStruct{
					Address:  "10.23.45.67",
					SshPort:  22,
					LabelIds: []request.ReqUint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.req.Address).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "fail3",
			s:    &s,
			args: args{
				req: &request.CreateNodeRequestStruct{
					Address:  "10.23.45.67",
					SshPort:  22,
					LabelIds: []request.ReqUint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.req.Address).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(1).
					WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail4",
			s:    &s,
			args: args{
				req: &request.CreateNodeRequestStruct{
					Address:  "10.23.45.67",
					SshPort:  22,
					LabelIds: []request.ReqUint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.req.Address).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_node`").WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				req: &request.CreateNodeRequestStruct{
					Address:  "10.23.45.67",
					SshPort:  22,
					LabelIds: []request.ReqUint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.req.Address).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_node`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if err := tt.s.CreateNode(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("CreateNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_DeleteNodeByIds(t *testing.T) {
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
			name: "fail1",
			s:    &s,
			args: args{ids: []uint{1}},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(1).
					WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail2",
			s:    &s,
			args: args{ids: []uint{1}},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "address"}).
					AddRow(1, "12.34.56.78"))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_node`").WithArgs(sqlmock.AnyArg(), tests2.AnyTime{}, 1).
					WillReturnError(errors.New("DB update error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "fail3",
			s:    &s,
			args: args{ids: []uint{1}},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "address"}).AddRow(1, "12.34.56.78"))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_node`").WithArgs(sqlmock.AnyArg(), tests2.AnyTime{}, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				// soft delete is update sql
				mock.ExpectExec("UPDATE `tb_sys_node`").WithArgs(tests2.AnyTime{}, 1).
					WillReturnError(errors.New("DB delete error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{ids: []uint{1}},
			invoke: func(args2 args) {
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if err := tt.s.DeleteNodeByIds(tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("DeleteNodeByIds() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_GetNodes(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		req *request.NodeListRequestStruct
	}
	n := uint(1)
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
				req: &request.NodeListRequestStruct{
					Address:     "12.34.56.78",
					Manager:     "tester",
					Region:      "Chengdu",
					Performance: &n,
					Asset:       "A1234",
					Creator:     "tester",
					Health:      &n,
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(
					fmt.Sprintf("%%%s%%", args2.req.Address),
					fmt.Sprintf("%%%s%%", args2.req.Manager),
					fmt.Sprintf("%%%s%%", args2.req.Region),
					args2.req.Performance,
					fmt.Sprintf("%%%s%%", args2.req.Asset),
					fmt.Sprintf("%%%s%%", args2.req.Creator),
					args2.req.Health,
				).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				req: &request.NodeListRequestStruct{
					Address:     "12.34.56.78",
					Manager:     "tester",
					Region:      "Chengdu",
					Performance: &n,
					Asset:       "A1234",
					Creator:     "tester",
					Health:      &n,
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(
					fmt.Sprintf("%%%s%%", args2.req.Address),
					fmt.Sprintf("%%%s%%", args2.req.Manager),
					fmt.Sprintf("%%%s%%", args2.req.Region),
					args2.req.Performance,
					fmt.Sprintf("%%%s%%", args2.req.Asset),
					fmt.Sprintf("%%%s%%", args2.req.Creator),
					args2.req.Health,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if _, err := tt.s.GetNodes(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("GetNodes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_RefreshNodeInfoById(t *testing.T) {
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
			name: "fail1",
			s:    &s,
			args: args{id: uint(1)},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.id).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{id: uint(1)},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.id).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tests2.SetLog()
			global.Machinery = nil
			global.Mysql = s.DB
			tt.invoke(tt.args)
			if err := s.RefreshNodeInfoById(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("RefreshNodeInfoById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// nolint: funlen
func TestMysqlService_UpdateNodeById(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		nodeId uint
		req    *request.UpdateNodeRequestStruct
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
				nodeId: uint(1),
				req: &request.UpdateNodeRequestStruct{
					LabelIds: []request.ReqUint{1},
					Manager:  "tester",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(args2.req.LabelIds[0]).
					WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail1",
			s:    &s,
			args: args{
				nodeId: uint(1),
				req: &request.UpdateNodeRequestStruct{
					LabelIds: []request.ReqUint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(args2.req.LabelIds[0]).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.nodeId).
					WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail2",
			s:    &s,
			args: args{
				nodeId: uint(1),
				req: &request.UpdateNodeRequestStruct{
					LabelIds: []request.ReqUint{1},
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(args2.req.LabelIds[0]).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.nodeId).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name: "fail3",
			s:    &s,
			args: args{
				nodeId: uint(1),
				req: &request.UpdateNodeRequestStruct{
					LabelIds: []request.ReqUint{1},
					Manager:  "tester",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(args2.req.LabelIds[0]).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.nodeId).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_node`").
					WithArgs(tests2.AnyTime{}, args2.req.LabelIds[0], args2.req.Manager).
					WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "fail4",
			s:    &s,
			args: args{
				nodeId: uint(1),
				req: &request.UpdateNodeRequestStruct{
					LabelIds: []request.ReqUint{1},
					Manager:  "tester",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(args2.req.LabelIds[0]).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.nodeId).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_node`").
					WithArgs(tests2.AnyTime{}, args2.req.LabelIds[0], args2.req.Manager).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_label`").WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				nodeId: uint(1),
				req: &request.UpdateNodeRequestStruct{
					LabelIds: []request.ReqUint{1},
					Manager:  "tester",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(args2.req.LabelIds[0]).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(args2.nodeId).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tb_sys_node`").
					WithArgs(tests2.AnyTime{}, args2.req.LabelIds[0], args2.req.Manager).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tb_sys_label`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if err := s.UpdateNodeById(tt.args.nodeId, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNodeById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
