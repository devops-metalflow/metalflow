package request

import "metalflow/pkg/response"

type CreateCronShutNodeRequest struct {
	Name      string   `json:"name" form:"name" validate:"required"`
	Keyword   string   `json:"keyword" form:"keyword" validate:"required"`
	StartTime string   `json:"startTime" form:"startTime" validate:"required"`
	ShutTime  string   `json:"shutTime" form:"shutTime" validate:"required"`
	NodeIds   []uint   `json:"nodeIds" form:"nodeIds" validate:"required"`
	Status    *ReqUint `json:"status" form:"status" validate:"required"`
	Creator   string   `json:"creator,omitempty" form:"creator"`
}

type ListCronShutNodeRequest struct {
	Name              string `json:"name" form:"name"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

type UpdateCronShutNodeRequest struct {
	Name      string   `json:"name" form:"name"`
	StartTime string   `json:"startTime" form:"startTime"`
	ShutTime  string   `json:"shutTime" form:"shutTime"`
	NodeIds   []uint   `json:"nodeIds" form:"nodeIds"`
	Status    *ReqUint `json:"status" form:"status"`
}

func (s *CreateCronShutNodeRequest) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "任务名称"
	m["Keyword"] = "任务名关键字"
	m["StartTime"] = "开机时间"
	m["ShutTime"] = "关机时间"
	m["Status"] = "任务状态"
	m["NodeIds"] = "机器节点"
	m["Date"] = "重复周期"
	return m
}
