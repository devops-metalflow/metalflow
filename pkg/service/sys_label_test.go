package service

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"metalflow/pkg/request"
	tests2 "metalflow/tests"
	"testing"
)

func TestMysqlService_GetLabels(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	type args struct {
		req *request.LabelListRequestStruct
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
			args: args{req: &request.LabelListRequestStruct{
				Name:    "test",
				Creator: "tester",
			}},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(
					fmt.Sprintf("%%%s%%", args2.req.Name),
					fmt.Sprintf("%%%s%%", args2.req.Creator),
				).WillReturnError(errors.New("DB search error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			args: args{
				req: &request.LabelListRequestStruct{
					Name:    "test",
					Creator: "tester",
				},
			},
			invoke: func(args2 args) {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_label`").WithArgs(
					fmt.Sprintf("%%%s%%", args2.req.Name),
					fmt.Sprintf("%%%s%%", args2.req.Creator),
				).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke(tt.args)
			if _, err := tt.s.GetLabels(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetLabels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
