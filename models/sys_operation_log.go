package models

import (
	"gorm.io/datatypes"
)

type SysOperationLog struct {
	Model
	ApiDesc   string         `json:"apiDesc,omitempty" gorm:"comment:'接口说明'"`
	Path      string         `json:"path,omitempty" gorm:"comment:'访问路径'"`
	Method    string         `json:"method,omitempty" gorm:"comment:'请求方式'"`
	Header    datatypes.JSON `json:"header,omitempty" gorm:"type:blob;comment:'请求header'"`
	Body      datatypes.JSON `json:"body,omitempty" gorm:"type:blob;comment:'请求主体'"`
	Data      datatypes.JSON `json:"data,omitempty" gorm:"type:blob;comment:'响应数据'"`
	Status    int            `json:"status,omitempty" gorm:"comment:'响应状态码'"`
	UserName  string         `json:"userName,omitempty" gorm:"comment:'用户名'"`
	RoleName  string         `json:"roleName,omitempty" gorm:"comment:'用户所属角色'"`
	Latency   int64          `json:"latency,omitempty" gorm:"comment:'请求耗时(ms)'"`
	UserAgent string         `json:"userAgent,omitempty" gorm:"comment:'浏览器标识'"`
}

func (m *SysOperationLog) TableName() string {
	return m.Model.TableName("sys_operation_log")
}
