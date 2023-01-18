package models

import "gorm.io/datatypes"

type SysNodeSecure struct {
	Model
	NodeId     uint           `gorm:"comment:'机器id'" json:"nodeId"`
	Node       SysNode        `gorm:"foreignkey:NodeId" json:"node"`
	Images     datatypes.JSON `gorm:"comment:'机器的docker镜像'" json:"images"`
	BareSbom   datatypes.JSON `gorm:"comment:'裸金属SBOM报告'" json:"bareSbom"`
	BareVul    datatypes.JSON `gorm:"comment:'裸金属vul报告'" json:"bareVul"`
	DockerSbom datatypes.JSON `gorm:"comment:'dockerSBOM报告'" json:"dockerSbom"`
	DockerVul  datatypes.JSON `gorm:"comment:'docker vul报告'" json:"dockerVul"`
	Score      uint           `gorm:"comment:'安全分数';type:tinyint(1);default:100" json:"score"`
}

func (m *SysNodeSecure) TableName() string {
	return m.Model.TableName("sys_node_secure")
}
