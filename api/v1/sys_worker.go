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

// GetWorkers get all workers.
func GetWorkers(c *gin.Context) {
	var req request.WorkerListRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	s := service.New(c)
	workers, err := s.GetWorkers(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// convert to ResponseStruct, hide some fields.
	var respStruct []response.WorkerListResponseStruct
	utils.Struct2StructByJson(workers, &respStruct)

	var resp response.PageData

	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// CreateWorker creates a worker.
func CreateWorker(c *gin.Context) {
	user := GetCurrentUser(c)
	// bind params.
	var req request.CreateWorkerRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	// params validate.
	err = global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// record current creator information.
	req.Creator = user.Username

	s := service.New(c)
	err = s.Create(req, new(models.SysWorker))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateWorkerById is used to update the corresponding worker information through workerId.
func UpdateWorkerById(c *gin.Context) {
	var req request.UpdateWorkerRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	// Get the workerId in the path.
	workerId := utils.Str2Uint(c.Param("workerId"))
	if workerId == 0 {
		response.FailWithMsg("the workerId is incorrect")
		return
	}

	s := service.New(c)
	// update data.
	err = s.UpdateById(workerId, req, new(models.SysWorker))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteWorkerByIds is used to delete workers in batches according to id.
func BatchDeleteWorkerByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	s := service.New(c)
	// delete data.
	err = s.DeleteByIds(req.GetUintIds(), new(models.SysWorker))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
