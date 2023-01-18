package v1

import (
	"github.com/gin-gonic/gin"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"
)

// CreateCronShutTask used to create a remote timer switch task.
func CreateCronShutTask(c *gin.Context) {
	// bind request body to struct.
	var req request.CreateCronShutNodeRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	u := GetCurrentUser(c)
	req.Creator = u.Username

	s := service.New(c)
	err = s.CreateCronShutNode(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// GetCronShutTasks get all remote timer switch tasks.
func GetCronShutTasks(c *gin.Context) {
	var req request.ListCronShutNodeRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	s := service.New(c)
	cronShutNode, err := s.GetCronShutNode(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// convert to ResponseStruct, hide some fields.
	var respStruct []response.CronShutNodeResponse
	utils.Struct2StructByJson(cronShutNode, &respStruct)
	// return paged data.
	var resp response.PageData
	// set pagination parameters.
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// UpdateCronShutTaskById used to update the remote timing switch task.
func UpdateCronShutTaskById(c *gin.Context) {
	var req request.UpdateCronShutNodeRequest
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	// get the shutId in the path.
	shutId := utils.Str2Uint(c.Param("shutId"))
	if shutId == 0 {
		response.FailWithMsg("obtain the timer switch task id error")
		return
	}

	s := service.New(c)
	// update data.
	err = s.UpdateCronShutNodeById(shutId, &req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteCronShutTask used to delete remote timing switch tasks in batches.
func BatchDeleteCronShutTask(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	s := service.New(c)
	// delete data.
	err = s.DeleteCronShutTaskByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
