package request

import (
	"metalflow/models"
	"metalflow/pkg/response"
)

// UserAuthRequestStruct 认证结构体
type UserAuthRequestStruct struct {
	Username string `json:"username" validate:"required"` // 用户
	Password string `json:"password" validate:"required"` // 密码
}

// UserListRequestStruct 获取用户列表结构体
type UserListRequestStruct struct {
	CurrentRole       models.SysRole
	Username          string `json:"username" form:"username"`
	Email             string `json:"email" form:"email"`
	OfficeName        string `json:"officeName" form:"officeName"`
	OrgName           string `json:"orgName" form:"orgName"`
	Introduction      string `json:"introduction" form:"introduction"`
	Status            *uint  `json:"status" form:"status"`
	RoleId            uint   `json:"roleId" form:"roleId"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

// CreateUserRequestStruct 创建用户结构体
type CreateUserRequestStruct struct {
	Username     string   `json:"username" validate:"required"`
	Password     string   `json:"password" validate:"required"`
	NickName     string   `json:"nickName"`
	Email        string   `json:"email"`
	OrgName      string   `json:"orgName"`
	OfficeName   string   `json:"officeName"`
	WorkPlace    string   `json:"workPlace"`
	Introduction string   `json:"introduction"`
	Status       *ReqUint `json:"status"`
	RoleId       uint     `json:"roleId" validate:"required"`
	Creator      string   `json:"creator"`
}

func (s *CreateUserRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Username"] = "用户名"
	m["Password"] = "初始密码"
	m["RoleId"] = "角色"
	return m
}

type UpdateUserRequestStruct struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	OfficeName string `json:"officeName"`
	OrgName    string `json:"orgName"`
	Status     *uint  `json:"status,omitempty"`
	RoleId     *uint  `json:"roleId,omitempty" validate:"required"`
	NickName   string `json:"nickName"`
}
