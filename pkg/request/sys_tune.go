package request

import "metalflow/pkg/response"

type TuneRollbackRequest struct {
	LogId uint `json:"logId" form:"logId"`
}

type TuneLogsRequest struct {
	response.PageInfo
}

type TuneSceneRequest struct {
	Scene string `json:"scene" form:"scene"`
}

type TuneAutoSetRequest struct {
	IsSave bool `json:"isSave" form:"isSave"`
}
