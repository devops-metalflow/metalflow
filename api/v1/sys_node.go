package v1

import (
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetNodes used to get all server nodes.
func GetNodes(c *gin.Context) {
	// bind request body to struct.
	var req request.NodeListRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	s := service.New(c)
	machines, err := s.GetNodes(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// convert to ResponseStruct, hide some fields.
	var respStruct []response.NodeListResponseStruct
	utils.Struct2StructByJson(machines, &respStruct)
	// return paged data.
	var resp response.PageData
	// set pagination parameters.
	resp.PageInfo = req.PageInfo
	// return response data.
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// CreateNode used to create a server node.
func CreateNode(c *gin.Context) {
	user := GetCurrentUser(c)
	// bind request body to struct.
	var req request.CreateNodeRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	// data validate.
	err = global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// record current creator information.
	req.Creator = user.Username

	s := service.New(c)
	err = s.CreateNode(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateNodeById update node info.
func UpdateNodeById(c *gin.Context) {
	// bind request body to struct.
	var req request.UpdateNodeRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	// get the nodeId in the path.
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}

	s := service.New(c)
	// update data.
	err = s.UpdateNodeById(nodeId, &req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteNodeByIds batch delete node by nodeId.
func BatchDeleteNodeByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("Parameter binding failed, please check the data type")
		return
	}

	s := service.New(c)
	// delete data.
	err = s.DeleteNodeByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// RefreshNodeInfo refresh node information.
func RefreshNodeInfo(c *gin.Context) {
	// get the nodeId in the path.
	nodeId := utils.Str2Uint(c.Param("nodeId"))
	if nodeId == 0 {
		response.FailWithMsg("the nodeId is incorrect")
		return
	}
	s := service.New(c)
	err := s.RefreshNodeInfoById(nodeId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchRebootNodeByIds batch reboot nodes by node id in database.
func BatchRebootNodeByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("参数绑定失败, 请检查数据类型")
		return
	}

	s := service.New(c)
	// delete node data.
	err = s.BatchRebootNodesByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
