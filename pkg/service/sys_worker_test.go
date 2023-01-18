package service

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"metalflow/pkg/request"
	tests2 "metalflow/tests"
	"testing"
)

func TestMysqlService_GetWorkers(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	n := uint(1)
	type args struct {
		req *request.WorkerListRequestStruct
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
			args: args{req: &request.WorkerListRequestStruct{
				Name:       "test",
				AutoDeploy: &n,
				Creator:    "tester",
			}},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_worker`").WithArgs(
					fmt.Sprintf("%%%s%%", args2.req.Name),
					args2.req.AutoDeploy,
					fmt.Sprintf("%%%s%%", args2.req.Creator),
				).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{req: &request.WorkerListRequestStruct{
				Name:       "test",
				AutoDeploy: &n,
				Creator:    "tester",
			}},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_worker`").WithArgs(
					fmt.Sprintf("%%%s%%", args2.req.Name),
					args2.req.AutoDeploy,
					fmt.Sprintf("%%%s%%", args2.req.Creator),
				).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if _, err := tt.s.GetWorkers(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("GetWorkers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
