// Package v1 contains all http handlers of the same version.
package v1

import (
	"github.com/gin-gonic/gin"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"
)

// GetApis gets api list.
func GetApis(c *gin.Context) {
	// bind request body to struct.
	var req request.ApiRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	// Create a data processing service.
	s := service.New(c)
	apis, err := s.GetApis(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// convert to ResponseStruct, hide some no used fields.
	var respStruct []response.ApiListResponseStruct
	utils.Struct2StructByJson(apis, &respStruct)
	// return paged data.
	var resp response.PageData
	// Set pagination parameters.
	resp.PageInfo = req.PageInfo
	// set data list, return response.
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// GetAllApiGroupByCategoryByRoleId queries the apis of the specified role (grouped by category).
func GetAllApiGroupByCategoryByRoleId(c *gin.Context) {
	// bind current user role sorting. hide specific users.
	user := GetCurrentUser(c)
	// Create a data processing service.
	s := service.New(c)
	// bind request body to struct.
	apis, ids, err := s.GetAllApiGroupByCategoryByRoleId(&user.Role, utils.Str2Uint(c.Param("roleId")))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	var resp response.ApiTreeWithAccessResponseStruct
	resp.AccessIds = ids
	utils.Struct2StructByJson(apis, &resp.List)
	response.SuccessWithData(resp)
}

// CreateApi for creating an api.
func CreateApi(c *gin.Context) {
	user := GetCurrentUser(c)
	// bind request body to struct.
	var req request.CreateApiRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	// parameter validator.
	err = global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// record current creator information.
	req.Creator = user.Username
	// Create a data processing service.
	s := service.New(c)
	err = s.CreateApi(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateApiById for updating an api.
func UpdateApiById(c *gin.Context) {
	// bind request body to struct.
	var req request.UpdateApiRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	// get the apiId in the path.
	apiId := utils.Str2Uint(c.Param("apiId"))
	if apiId == 0 {
		response.FailWithMsg("the api id is incorrect")
		return
	}
	// Create a data processing service.
	s := service.New(c)
	// update data.
	err = s.UpdateApiById(apiId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteApiByIds used to delete apis in batches.
func BatchDeleteApiByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type.")
		return
	}

	// Create a data processing service.
	s := service.New(c)
	// delete data.
	err = s.DeleteApiByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
