package models //nolint:gofmt

// SysCollection collect user favorite nodes
type SysCollection struct {
	Model
	Username    string `gorm:"comment:'工号'" json:"username"`
	NodeId      uint   `gorm:"comment:'机器id'" json:"nodeId"`
	Description string `gorm:"comment:'机器描述'" json:"description"`
}

func (m *SysCollection) TableName() string {
	return m.Model.TableName("sys_collection")
}
