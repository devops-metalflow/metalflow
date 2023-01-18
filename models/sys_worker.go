package models

import (
	"fmt"
	"metalflow/pkg/global"
)

type SysWorker struct {
	Model
	Name       string     `json:"name" gorm:"comment:'worker名称';unique"`
	Desc       string     `json:"desc" gorm:"comment:'worker描述'"`
	Port       int        `json:"port" gorm:"comment:'worker暴露端口'"`
	AutoDeploy *uint      `json:"autoDeploy" gorm:"type:tinyint(1);default:1;comment:'是否节点注册时部署'"`
	DeployCmd  string     `json:"deployCmd" gorm:"comment:'部署命令'"`
	StartCmd   string     `json:"startCmd" gorm:"comment:'启动命令'"`
	StopCmd    string     `json:"stopCmd" gorm:"comment:'停止命令'"`
	ReloadCmd  string     `json:"reloadCmd" gorm:"comment:'重新加载命令'"`
	DeleteCmd  string     `json:"deleteCmd" gorm:"comment:'删除命令'"`
	ServiceReq string     `json:"serviceReq" gorm:"comment:'grpc请求体'"`
	CheckReq   string     `json:"checkReq" gorm:"comment:'worker状态检查请求体'"`
	Creator    string     `json:"creator" gorm:"comment:'创建人'"`
	Nodes      []*SysNode `json:"nodes" gorm:"many2many:sys_node_worker_relation"`
}

func (m *SysWorker) TableName() string {
	return m.Model.TableName("sys_worker")
}

// RelationNodeWorker save node id which had deployed workers normally.
type RelationNodeWorker struct {
	SysWorkerId uint `json:"sysWorkerId"`
	SysNodeId   uint `json:"sysNodeId"`
}

func (m RelationNodeWorker) TableName() string {
	return fmt.Sprintf("%s_%s", global.Conf.Mysql.TablePrefix, "sys_node_worker_relation")
}
