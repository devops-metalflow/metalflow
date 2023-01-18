package service

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	tests2 "metalflow/tests"
	"testing"
)

func TestMysqlService_GetCountData(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	tests := []struct {
		name    string
		s       *MysqlService
		invoke  func()
		wantErr bool
	}{
		{
			name: "fail1",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail2",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail3",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(0).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "fail4",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(0).WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			if _, err := tt.s.GetCountData(); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetCountData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_GetHealthNodeCount(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	tests := []struct {
		name    string
		s       *MysqlService
		invoke  func()
		wantErr bool
	}{
		{
			name: "fail1",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			if _, err := tt.s.GetHealthNodeCount(); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetHealthNodeCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_GetManagerNodeCount(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	tests := []struct {
		name    string
		s       *MysqlService
		invoke  func()
		wantErr bool
	}{
		{
			name: "fail1",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			if _, err := tt.s.GetManagerNodeCount(); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetManagerNodeCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_GetPerformanceNodeCount(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	tests := []struct {
		name    string
		s       *MysqlService
		invoke  func()
		wantErr bool
	}{
		{
			name: "fail1",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			if _, err := tt.s.GetPerformanceNodeCount(); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetPerformanceNodeCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlService_GetRegionNodeCount(t *testing.T) {
	mock := tests2.GetMock()
	s := New(nil)
	tests := []struct {
		name    string
		s       *MysqlService
		invoke  func()
		wantErr bool
	}{
		{
			name: "fail1",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnError(errors.New("DB error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			s:    &s,
			invoke: func() {
				mock.ExpectQuery("SELECT (.*) FROM `tb_sys_node`").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.invoke()
			if _, err := tt.s.GetRegionNodeCount(); (err != nil) != tt.wantErr {
				t.Errorf("MysqlService.GetRegionNodeCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
