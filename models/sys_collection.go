package models

// SysCollection collect user favorite nodes
type SysCollection struct {
	Model
	Username string `gorm:"comment:'工号'" json:"username"`
	NodeId   uint   `gorm:"comment:'机器id'" json:"nodeId"`
}

func (m *SysCollection) TableName() string {
	return m.Model.TableName("sys_collection")
}
