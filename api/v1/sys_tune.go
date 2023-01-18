package v1

import (
	"github.com/gin-gonic/gin"
	"metalflow/models"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"
)

// GetTuneScoreByNodeId obtain the tuning score of the corresponding node according to nodeId.
func GetTuneScoreByNodeId(c *gin.Context) {
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	score, err := s.GetTuneScoreByNodeId(nodeId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(score)
}

// GetTuneLogsByNodeId obtain the tuning record of the corresponding machine.
func GetTuneLogsByNodeId(c *gin.Context) {
	var req request.TuneLogsRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	list, err := s.GetTuneLogsByNodeId(nodeId, &req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// convert to ResponseStruct, hide some fields.
	var respStruct []response.TuneLogList
	utils.Struct2StructByJson(list, &respStruct)
	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// Cleanup garbage the node.
func Cleanup(c *gin.Context) {
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	err := s.Cleanup(nodeId)
	if err != nil {
		response.FailWithMsg("cleanup failed：" + err.Error())
		return
	}
	response.Success()
}

// Rollback roll back tuning operations.
func Rollback(c *gin.Context) {
	var req request.TuneRollbackRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	err = s.Rollback(nodeId, req.LogId)
	if err != nil {
		response.FailWithMsg("rollback failed：" + err.Error())
		return
	}
	response.Success()
}

// Set intelligent tuning of nodes.
func Set(c *gin.Context) {
	var req request.TuneAutoSetRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	err = s.SetTune(nodeId, &req)
	if err != nil {
		response.FailWithMsg("intelligent tuning failed：" + err.Error())
		return
	}
	response.SuccessWithData("in tuning, please check the results later...")
}

// Scene customize the corresponding scene for the node.
func Scene(c *gin.Context) {
	var req request.TuneSceneRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	err = s.SetScene(nodeId, req.Scene)
	if err != nil {
		response.FailWithMsg("scene tuning failed：" + err.Error())
		return
	}
	response.Success()
}

// Turbo accelerate the performance of nodes.
func Turbo(c *gin.Context) {
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	err := s.Turbo(nodeId)
	if err != nil {
		response.FailWithMsg("performance acceleration failed：" + err.Error())
		return
	}
	response.Success()
}

// BatchDeleteTuneLogByIds batch delete tuning records.
func BatchDeleteTuneLogByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	s := service.New(c)
	// delete data.
	err = s.DeleteByIds(req.GetUintIds(), new(models.SysNodeTuneLog))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
