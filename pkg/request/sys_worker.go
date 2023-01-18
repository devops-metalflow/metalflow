package request

import "metalflow/pkg/response"

type WorkerListRequestStruct struct {
	Name              string `json:"name,omitempty" form:"name"`
	AutoDeploy        *uint  `json:"autoDeploy,omitempty" form:"autoDeploy"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

type CreateWorkerRequestStruct struct {
	Name       string `json:"name" form:"name" validate:"required"`
	Desc       string `json:"desc" form:"desc"`
	AutoDeploy *uint  `json:"autoDeploy" form:"autoDeploy" validate:"required"`
	Port       int    `json:"port" form:"port" validate:"required"`
	DeployCmd  string `json:"deployCmd" form:"deployCmd" validate:"required"`
	StartCmd   string `json:"startCmd" form:"startCmd" validate:"required"`
	StopCmd    string `json:"stopCmd" form:"stopCmd" validate:"required"`
	DeleteCmd  string `json:"deleteCmd" form:"deleteCmd"`
	ServiceReq string `json:"serviceReq" form:"serviceReq"`
	CheckReq   string `json:"checkReq" form:"checkReq"`
	Creator    string `json:"creator" form:"creator"`
}

type UpdateWorkerRequestStruct struct {
	Name       string `json:"name" form:"name"`
	Desc       string `json:"desc" form:"desc"`
	AutoDeploy *uint  `json:"autoDeploy" form:"autoDeploy"`
	Port       int    `json:"port" form:"port"`
	DeployCmd  string `json:"deployCmd" form:"deployCmd"`
	StartCmd   string `json:"startCmd" form:"startCmd"`
	StopCmd    string `json:"stopCmd" form:"stopCmd"`
	DeleteCmd  string `json:"deleteCmd" form:"deleteCmd"`
	ServiceReq string `json:"serviceReq" form:"serviceReq"`
	CheckReq   string `json:"checkReq" form:"checkReq"`
}

func (s *CreateWorkerRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "worker名称"
	m["Port"] = "worker服务端口"
	m["DeployCmd"] = "worker部署命令"
	m["StartCmd"] = "worker启动命令"
	m["StopCmd"] = "worker停止命令"
	return m
}
