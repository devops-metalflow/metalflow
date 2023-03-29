package request

import (
	"metalflow/pkg/response"
)

// NodeListRequestStruct 获取机器列表结构体
type NodeListRequestStruct struct {
	Address           string `json:"address" form:"address"`
	SshPort           int    `json:"sshPort" form:"sshPort"`
	Os                string `json:"os" form:"os"`
	Asset             string `json:"asset" form:"asset"`
	Health            *uint  `json:"health" form:"health"`
	Manager           string `json:"manager" form:"manager"`
	Performance       *uint  `json:"performance" form:"performance"`
	Region            string `json:"region" form:"region"`
	Remark            string `json:"remark" form:"remark"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

// CreateNodeRequestStruct 创建机器结构体
type CreateNodeRequestStruct struct {
	Address     string    `json:"address" form:"address" validate:"required"`
	SshPort     uint      `json:"sshPort" form:"sshPort" validate:"required"`
	Os          string    `json:"os" form:"os" validate:"required"`
	ServicePort int       `json:"servicePort" form:"servicePort"`
	Asset       string    `json:"asset" form:"asset"`
	Health      *uint     `json:"health" form:"health"`
	Performance *uint     `json:"performance" form:"performance"`
	Region      string    `json:"region" form:"region"`
	Remark      string    `json:"remark" form:"remark"`
	Creator     string    `json:"creator" form:"creator"`
	LabelIds    []ReqUint `json:"labelIds" form:"labelIds"`
}

type NodeShellConnectRequestStruct struct {
	Address  string  `json:"address" form:"address" validate:"required"`
	SshPort  ReqUint `json:"sshPort" form:"sshPort" validate:"required"`
	Username string  `json:"username" form:"username" validate:"required"`
	Password string  `json:"password" form:"password" validate:"required"`
}

// NodeShellWsRequestStruct 机器shell_ws请求结构体
type NodeShellWsRequestStruct struct {
	Address string  `json:"address" form:"address"`
	Cols    ReqUint `json:"cols,omitempty" form:"cols"`
	Rows    ReqUint `json:"rows,omitempty" form:"rows"`
	SshId   string  `json:"sshId" form:"sshId"`
}

type NodeShellFileStruct struct {
	Path  string `json:"path" form:"path"`
	SshId string `json:"sshId" form:"sshId"`
}

type ResizeWsStruct struct {
	High  ReqUint `json:"high,omitempty" form:"high"`
	Width ReqUint `json:"width,omitempty" form:"high"`
	SshId string  `json:"sshId" form:"sshId"`
}

type ModifyFileRequestStruct struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	SshId   string `json:"sshId" form:"sshId"`
}

type NodeVncWsRequestStruct struct {
	Address string `json:"address" form:"address" validate:"required"`
	Port    int    `json:"port" form:"port" validate:"required"`
}

// UpdateNodeRequestStruct 更新机器结构体
type UpdateNodeRequestStruct struct {
	LabelIds []ReqUint `json:"labelIds"` // 标签的ids
	Manager  string    `json:"manager"`  // 机器对应责任人
}

// FieldTrans 翻译需要校验的字段名称
func (s *CreateNodeRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Address"] = "主机ip"
	m["SshPort"] = "ssh端口"
	m["Os"] = "操作系统"
	return m
}
