package v1

import (
	"github.com/gin-gonic/gin"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"
)

// GetOperationLogs used to get operation log.
func GetOperationLogs(c *gin.Context) {
	// bind params
	var req request.OperationLogRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("Parameter binding failed, please check the data type")
		return
	}

	s := service.New(c)
	operationLogs, err := s.GetOperationLogs(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	var respStruct []response.OperationLogListResponseStruct
	utils.Struct2StructByJson(operationLogs, &respStruct)

	var resp response.PageData

	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// BatchDeleteOperationLogByIds batch delete operation logs
func BatchDeleteOperationLogByIds(c *gin.Context) {
	if !global.Conf.System.OperationLogAllowedToDelete {
		response.FailWithMsg("log deletion has been turned off by the administrator")
		return
	}
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	s := service.New(c)
	// delete data.
	err = s.DeleteByIds(req.GetUintIds(), new(models.SysOperationLog))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
