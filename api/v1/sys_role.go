package v1

import (
	"fmt"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetRoles get a list of roles.
func GetRoles(c *gin.Context) {
	// bind params.
	var req request.RoleListRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("Parameter binding failed, please check the data type")
		return
	}

	// bind current user role sorting. Hide specific users.
	user := GetCurrentUser(c)
	req.CurrentRoleSort = *user.Role.Sort

	// create service.
	s := service.New(c)
	roles, err := s.GetRoles(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	var respStruct []response.RoleListResponseStruct
	utils.Struct2StructByJson(roles, &respStruct)

	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// CreateRole is used to create a role.
func CreateRole(c *gin.Context) {
	user := GetCurrentUser(c)
	// bind params.
	var req request.CreateRoleRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("Bind current user role sorting. Hide specific users")
		return
	}

	// param validate.
	err = global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	if req.Sort != nil && *user.Role.Sort > uint(*req.Sort) {
		response.FailWithMsg(fmt.Sprintf("role sorting is not allowed to be smaller than the current login account number (%d)", *user.Role.Sort))
		return
	}

	// record current creator information.
	req.Creator = user.Username

	s := service.New(c)
	err = s.Create(req, new(models.SysRole))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateRoleById update role.
func UpdateRoleById(c *gin.Context) {
	var req request.UpdateRoleRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	if req.Sort != nil {
		// bind current user role sorting. Hide specific users.
		user := GetCurrentUser(c)
		if req.Sort != nil && *user.Role.Sort > uint(*req.Sort) {
			response.FailWithMsg(fmt.Sprintf("role sorting is not allowed to be smaller than current user sorting (%d)", *user.Role.Sort))
			return
		}
	}

	// get the roleId in the path.
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.FailWithMsg("incorrect roleId")
		return
	}

	user := GetCurrentUser(c)
	if req.Status != nil && uint(*req.Status) == models.SysRoleStatusDisabled && roleId == user.RoleId {
		response.FailWithMsg("Cannot disable own character")
		return
	}

	s := service.New(c)
	// update data.
	err = s.UpdateById(roleId, req, new(models.SysRole))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateRoleMenusById update a role's permissions menu.
func UpdateRoleMenusById(c *gin.Context) {
	// bind params.
	var req request.UpdateIncrementalIdsRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("参数绑定失败, %v", err))
		return
	}
	// get the roleId in the path.
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.FailWithMsg("incorrect character number")
		return
	}
	// bind current user role sorting. hide specific users.
	user := GetCurrentUser(c)

	if user.RoleId == roleId {
		if *user.Role.Sort == models.SysRoleSuperAdminSort && len(req.Delete) > 0 {
			response.FailWithMsg("unable to remove super administrator privileges, if you have any questions, please contact the website developer")
			return
		} else if *user.Role.Sort != models.SysRoleSuperAdminSort {
			response.FailWithMsg("unable to change your own permissions, if you need to change, please contact the superior leader")
			return
		}
	}

	s := service.New(c)
	err = s.UpdateRoleMenusById(&user.Role, roleId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateRoleApisById update the permission interface of the role.
func UpdateRoleApisById(c *gin.Context) {
	// bind request body.
	var req request.UpdateIncrementalIdsRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("parameter binding failed: %v", err))
		return
	}
	// get the roleId in the path.
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.FailWithMsg("incorrect character number")
		return
	}

	// bind current user role sorting. hide specific users.
	user := GetCurrentUser(c)

	if user.RoleId == roleId {
		if *user.Role.Sort == models.SysRoleSuperAdminSort && len(req.Delete) > 0 {
			response.FailWithMsg("unable to remove super administrator privileges, if you have any questions, please contact the website developer")
			return
		} else if *user.Role.Sort != models.SysRoleSuperAdminSort {
			response.FailWithMsg("unable to change your own permissions, if you need to change, please contact the superior leader")
			return
		}
	}

	s := service.New(c)
	err = s.UpdateRoleApisById(roleId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteRoleByIds batch delete roles.
func BatchDeleteRoleByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	user := GetCurrentUser(c)
	if utils.ContainsUint(req.GetUintIds(), user.RoleId) {
		response.FailWithMsg("cannot delete own role")
		return
	}

	// the dev role is the default role and cannot be deleted. the dev role id is 2.
	if utils.ContainsUint(req.GetUintIds(), global.DefaultDevRoleId) {
		response.FailWithMsg("The role of dev personnel is the default role and cannot be deleted=")
		return
	}

	s := service.New(c)
	// delete data.
	err = s.DeleteRoleByIds(req.GetUintIds())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
