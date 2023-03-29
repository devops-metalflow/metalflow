package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"metalflow/pkg/global"
)

type SysWorker struct {
	Model
	Name       string         `json:"name" gorm:"comment:'worker名称';unique"`
	Desc       string         `json:"desc" gorm:"comment:'worker描述'"`
	Port       int            `json:"port" gorm:"comment:'worker暴露端口'"`
	AutoDeploy *uint          `json:"autoDeploy" gorm:"type:tinyint(1);default:0;comment:'是否节点注册时部署'"`
	DeployCmd  datatypes.JSON `json:"deployCmd" gorm:"comment:'部署命令'"`
	ServiceReq string         `json:"serviceReq" gorm:"comment:'grpc请求体'"`
	CheckReq   string         `json:"checkReq" gorm:"comment:'worker状态检查请求体'"`
	Creator    string         `json:"creator" gorm:"comment:'创建人'"`
	Nodes      []*SysNode     `json:"nodes" gorm:"many2many:sys_node_worker_relation"`
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

type CmdStruct struct {
	Download string `json:"download,omitempty"`
	Start    string `json:"start,omitempty"`
	Stop     string `json:"stop,omitempty"`
}

func (m *SysWorker) GetCmd(os string) (*CmdStruct, error) {
	var osDeployCmd map[string]*CmdStruct
	err := json.Unmarshal([]byte(m.DeployCmd.String()), &osDeployCmd)
	if err != nil {
		return nil, err
	}
	cmd := osDeployCmd[os]
	return cmd, nil
}
