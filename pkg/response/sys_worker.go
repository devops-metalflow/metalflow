package response

type WorkerListResponseStruct struct {
	Id         uint   `json:"id" form:"id"`
	Name       string `json:"name" form:"name"`
	Desc       string `json:"desc" form:"desc"`
	Port       int    `json:"port" form:"port"`
	StartCmd   string `json:"startCmd" form:"startCmd"`
	DeployCmd  string `json:"deployCmd" form:"deployCmd"`
	StopCmd    string `json:"stopCmd" form:"stopCmd"`
	DeleteCmd  string `json:"deleteCmd" form:"deleteCmd"`
	AutoDeploy *uint  `json:"autoDeploy" form:"autoDeploy"`
	ServiceReq string `json:"serviceReq" form:"serviceReq"`
	CheckReq   string `json:"checkReq" form:"checkReq"`
	Creator    string `json:"creator" form:"creator"`
	CreatedAt  string `json:"createdAt" form:"createdAt"`
}
