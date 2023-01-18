package response

import "metalflow/models"

// UserInfoResponseStruct 用户信息响应
type UserInfoResponseStruct struct {
	Id           uint     `json:"id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	OfficeName   string   `json:"officeName"`
	OrgName      string   `json:"orgName"`
	WorkPlace    string   `json:"workPlace"`
	Introduction string   `json:"introduction"`
	NickName     string   `json:"nickName"`
	Roles        []string `json:"roles"`
	RoleSort     uint     `json:"roleSort"`
}

// UserListResponseStruct 用户信息响应, 字段含义见models.SysUser
type UserListResponseStruct struct {
	Id           uint             `json:"id"`
	Username     string           `json:"username"`
	Email        string           `json:"email"`
	OrgName      string           `json:"orgName"`
	OfficeName   string           `json:"officeName"`
	Introduction string           `json:"introduction"`
	WorkPlace    string           `json:"workPlace"`
	Status       *uint            `json:"status"`
	RoleId       uint             `json:"roleId"`
	Creator      string           `json:"creator"`
	NickName     string           `json:"nickName"`
	CreatedAt    models.LocalTime `json:"createdAt"`
}
