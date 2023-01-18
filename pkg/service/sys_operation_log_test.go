package service

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"metalflow/pkg/request"
	tests2 "metalflow/tests"
	"testing"
)

func TestMysqlService_GetOperationLogs(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		req *request.OperationLogRequestStruct
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
			args: args{req: &request.OperationLogRequestStruct{
				Method:   "POST",
				Path:     "v1",
				Username: "12345678",
				Status:   "200",
			}},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_operation_log`").WithArgs(
					fmt.Sprintf("%%%s%%", args2.req.Method),
					fmt.Sprintf("%%%s%%", args2.req.Username),
					fmt.Sprintf("%%%s%%", args2.req.Path),
					fmt.Sprintf("%%%s%%", args2.req.Status),
				).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{req: &request.OperationLogRequestStruct{
				Method:   "POST",
				Path:     "v1",
				Username: "12345678",
				Status:   "200",
			}},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_operation_log`").WithArgs(
					fmt.Sprintf("%%%s%%", args2.req.Method),
					fmt.Sprintf("%%%s%%", args2.req.Username),
					fmt.Sprintf("%%%s%%", args2.req.Path),
					fmt.Sprintf("%%%s%%", args2.req.Status),
				).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if _, err := tt.s.GetOperationLogs(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("GetOperationLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
