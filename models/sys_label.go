package models

import (
	"fmt"
	"metalflow/pkg/global"
)

type SysLabel struct {
	Model
	Name    string    `gorm:"comment:'标签名称'" json:"name,omitempty"`
	Creator string    `gorm:"comment:'创建人'" json:"creator"`
	Nodes   []SysNode `gorm:"many2many:sys_node_label_relation" json:"nodes,omitempty"`
}

func (l *SysLabel) TableName() string {
	return l.Model.TableName("sys_label")
}

type RelationNodeLabel struct {
	SysNodeId  uint `json:"sysNodeId,omitempty"`
	SysLabelId uint `json:"sysLabelId,omitempty"`
}

func (r RelationNodeLabel) TableName() string {
	return fmt.Sprintf("%s_%s", global.Conf.Mysql.TablePrefix, "sys_node_label_relation")
}
