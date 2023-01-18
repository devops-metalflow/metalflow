package models

const (
	// SysUserStatusDisabled 用户状态
	SysUserStatusDisabled    uint   = 0    // 禁用
	SysUserStatusNormal      uint   = 1    // 正常
	SysUserStatusDisabledStr string = "禁用" // 禁用
	SysUserStatusNormalStr   string = "正常" // 正常
)

// SysUser User
type SysUser struct {
	Model
	Username     string  `gorm:"unique;comment:'user login name'" json:"username"`
	Password     string  `gorm:"unique;comment:'password'" json:"password"`
	NickName     string  `gorm:"comment:'昵称'" json:"nickName"`
	Email        string  `gorm:"comment:'邮箱'" json:"email"`
	OfficeName   string  `gorm:"comment:'科室'" json:"officeName"`
	OrgName      string  `gorm:"comment:'部门'" json:"orgName"`
	WorkPlace    string  `gorm:"comment:'常驻地'" json:"workPlace"`
	Introduction string  `gorm:"comment:'自我介绍'" json:"introduction"`
	Status       *uint   `gorm:"default:1;comment:'用户状态(正常/禁用, 默认正常)'" json:"status"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	Creator      string  `gorm:"comment:'创建人'" json:"creator"`
	RoleId       uint    `gorm:"comment:'角色Id外键'" json:"roleId"`
	Role         SysRole `gorm:"foreignkey:RoleId" json:"role"` // 用户属于某个角色，将SysUser.RoleId指定为外键
}

func (m *SysUser) TableName() string {
	return m.Model.TableName("sys_user")
}
