package models

const (
	// SysRoleStatusDisabled 用户状态
	SysRoleStatusDisabled    uint   = 0    // 禁用
	SysRoleStatusNormal      uint   = 1    // 正常
	SysRoleStatusDisabledStr string = "禁用" // 禁用
	SysRoleStatusNormalStr   string = "正常" // 正常

	SysRoleSuperAdminSort uint = 0 // 超级管理员排序
)

// SysRole 系统角色表
type SysRole struct {
	Model
	Name    string     `gorm:"comment:'角色名称'" json:"name"`
	Keyword string     `gorm:"unique;comment:'角色关键词'" json:"keyword"`
	Desc    string     `gorm:"comment:'角色说明'" json:"desc"`
	Status  *uint      `gorm:"type:tinyint(1);default:1;comment:'角色状态(正常/禁用, 默认正常)'" json:"status"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	Sort    *uint      `gorm:"default:1;comment:'角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为0表示超级管理员)'" json:"sort"`
	Creator string     `gorm:"comment:'创建人'" json:"creator"`
	Menus   []*SysMenu `gorm:"many2many:sys_role_menu_relation;" json:"menus"` // 角色菜单多对多关系
	Users   []*SysUser `gorm:"foreignkey:RoleId"`                              // 一个角色有多个user
}

func (m *SysRole) TableName() string {
	return m.Model.TableName("sys_role")
}
