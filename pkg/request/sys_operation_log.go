package request

import (
	"metalflow/pkg/response"
	"time"
)

// OperationLogRequestStruct get operation log struct.
type OperationLogRequestStruct struct {
	Method            string `json:"method" form:"method"`
	Path              string `json:"path" form:"path"`
	Username          string `json:"username" form:"username"`
	Status            string `json:"status" form:"status"`
	response.PageInfo        // 分页参数
}

// FieldTrans for translate needed field.
func (s *OperationLogRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Status"] = "响应状态码"
	return m
}

// CreateOperationLogRequestStruct for create operation log.
type CreateOperationLogRequestStruct struct {
	ApiDesc   string        `json:"apiDesc"`
	Path      string        `json:"path"`
	Method    string        `json:"method"`
	Params    string        `json:"params"`
	Body      string        `json:"body"`
	Data      string        `json:"data"`
	Status    ReqUint       `json:"status"`
	Username  string        `json:"username"`
	RoleName  string        `json:"roleName"`
	Latency   time.Duration `json:"latency"`
	UserAgent string        `json:"userAgent"`
}
