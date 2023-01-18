package models

const (
	AutoTuneType = 0
	CleanupType  = 1
	TurboType    = 2
)

type SysNodeTuneLog struct {
	Model
	NodeId      uint    `gorm:"comment:'机器节点id'" json:"nodeId"`
	Node        SysNode `gorm:"foreignkey:NodeId" json:"node"`
	TuneType    uint    `gorm:"tinyint(1);comment:'请求类型'" json:"tuneType"`
	RespProfile string  `gorm:"comment:'处理后的profile'" json:"respProfile"`
}

func (m *SysNodeTuneLog) TableName() string {
	return m.Model.TableName("sys_node_tune_log")
}

type SysNodeTuneScene struct {
	Model
	Scene   string `gorm:"comment:'场景'" json:"scene"`
	Profile string `gorm:"comment:'场景对应的profile'" json:"profile"`
	Creator string `gorm:"comment:'创建人'" json:"creator"`
}

func (m *SysNodeTuneScene) TableName() string {
	return m.Model.TableName("sys_node_tune_scene")
}
