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

// GetLabels gets the list of labels.
func GetLabels(c *gin.Context) {
	// bind request body to struct.
	var req request.LabelListRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	s := service.New(c)
	labels, err := s.GetLabels(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = labels
	response.SuccessWithData(resp)
}

// CreateLabel creates machine labels.
func CreateLabel(c *gin.Context) {
	user := GetCurrentUser(c)
	// bind request body to struct.
	var req request.CreateLabelRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	// params validate.
	err = global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	req.Creator = user.Username
	s := service.New(c)
	err = s.Create(req, new(models.SysLabel))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateLabelById updates label name.
func UpdateLabelById(c *gin.Context) {
	var req request.UpdateLabelRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	labelId := utils.Str2Uint(c.Param("labelId"))
	s := service.New(c)
	err = s.UpdateById(labelId, req, new(models.SysLabel))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteLabelByIds used to delete tags in batch.
func BatchDeleteLabelByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	s := service.New(c)
	// delete data.
	err = s.DeleteByIds(req.GetUintIds(), new(models.SysLabel))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
