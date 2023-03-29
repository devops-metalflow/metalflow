package models

import "gorm.io/datatypes"

const (
	SysNodeHealthNormal   uint = 0 // 运行中
	SysNodeHealthAbnormal uint = 1 // 异常
	SysNodeHealthShutdown uint = 2 // 已停机
)

// SysNode 机器节点信息
type SysNode struct {
	Model
	Address         string         `gorm:"unique;comment:'主机地址(ip)'" json:"address"`
	Os              string         `gorm:"comment:'操作系统'" json:"os"`
	SshPort         uint           `gorm:"comment:'ssh端口号';default:22" json:"sshPort"`
	ServicePort     int            `gorm:"comment:'注册的服务的端口';default:19090" json:"servicePort"`
	Asset           string         `gorm:"comment:'资产编号'" json:"asset"`
	Manager         string         `gorm:"comment:'责任人'" json:"manager"`
	Health          *uint          `gorm:"type:tinyint(1);comment:'健康度(0:运行中 1:异常 2:已停机)';default:0" json:"health"`
	Performance     *uint          `gorm:"type:tinyint(1);comment:'性能(0:高 1:中 2:低)';default:0" json:"performance"`
	PingStat        *uint          `gorm:"type:tinyint(1);comment:'ping状态'" json:"pingStat"`
	Region          string         `gorm:"comment:'地域'" json:"region"`
	Remark          string         `gorm:"comment:'说明'" json:"remark"`
	Creator         string         `gorm:"comment:'创建人'" json:"creator"`
	Metrics         string         `gorm:"comment:'机器配置'" json:"metrics"`
	Information     datatypes.JSON `gorm:"comment:'机器详情'" json:"information"`
	Labels          []SysLabel     `gorm:"many2many:sys_node_label_relation" json:"labels"`
	RefreshLastTime LocalTime      `gorm:"comment:'上次刷新时间'" json:"refreshLastTime"`
	RefreshCount    *uint          `gorm:"comment:'刷新次数';default:0" json:"refreshCount"`
	Workers         []*SysWorker   `gorm:"many2many:sys_node_worker_relation" json:"workers"`
}

func (m *SysNode) TableName() string {
	return m.Model.TableName("sys_node")
}
