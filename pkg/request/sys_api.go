package request

import "metalflow/pkg/response"

// ApiRequestStruct 获取接口列表结构体
type ApiRequestStruct struct {
	Method            string `json:"method" form:"method"`
	Path              string `json:"path" form:"path"`
	Category          string `json:"category" form:"category"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

// CreateApiRequestStruct 创建接口结构体
type CreateApiRequestStruct struct {
	Method   string `json:"method" validate:"required"`
	Path     string `json:"path" validate:"required"`
	Category string `json:"category" validate:"required"`
	Desc     string `json:"desc"`
	Title    string `json:"title"`
	Creator  string `json:"creator"`
	RoleIds  []uint `json:"roleIds"` // 绑定可以访问该接口的角色
}

// FieldTrans 翻译需要校验的字段名称
func (s *CreateApiRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Method"] = "请求方式"
	m["Path"] = "访问路径"
	m["Category"] = "所属类别"
	return m
}

// UpdateApiRequestStruct 更新接口结构体
type UpdateApiRequestStruct struct {
	Method   *string `json:"method"`
	Path     *string `json:"path"`
	Category *string `json:"category"`
	Desc     *string `json:"desc"`
	Title    *string `json:"title"`
	Creator  *string `json:"creator"`
	RoleIds  []uint  `json:"roleIds"`
}
