package models

import (
	"fmt"
	"metalflow/pkg/global"
)

const (
	SysCronShutNodeDisable uint = 0
	SysCronShutNodeEnable  uint = 1
)

// SysCronShutNode 定时开关机任务，简单起见，这里直接用一张现成表，且不区分用户定制
type SysCronShutNode struct {
	Model
	Name      string     `json:"name" gorm:"comment:'任务名称'"`
	Keyword   string     `json:"keyword" gorm:"unique;comment:'任务名关键字，必须为英文且唯一'"`
	StartTime string     `json:"startTime" gorm:"comment:'开机时间'"`
	ShutTime  string     `json:"shutTime" gorm:"comment:'关机时间'"`
	Status    *uint      ` json:"status" gorm:"type:tinyint(1);default:1;comment:'任务状态(正常/禁用, 默认正常)'"`
	Creator   string     `json:"creator" gorm:"comment:'创建人'"`
	Nodes     []*SysNode `json:"nodes" gorm:"many2many:sys_node_shut_relation"`
}

func (m *SysCronShutNode) TableName() string {
	return m.Model.TableName("sys_cron_shut_node")
}

type RelationNodeShut struct {
	SysShutId uint `json:"sysShutId"`
	SysNodeId uint `json:"sysNodeId"`
}

func (m RelationNodeShut) TableName() string {
	return fmt.Sprintf("%s_%s", global.Conf.Mysql.TablePrefix, "sys_node_shut_relation")
}
