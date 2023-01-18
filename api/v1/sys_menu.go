package v1

import (
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetMenuTree queries the current user's menu tree.
func GetMenuTree(c *gin.Context) {
	user := GetCurrentUser(c)

	s := service.New(c)
	menus, err := s.GetMenuTree(user.RoleId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// covert to MenuTreeResponseStruct.
	var resp []response.MenuTreeResponseStruct
	utils.Struct2StructByJson(menus, &resp)
	response.SuccessWithData(resp)
}

// GetAllMenuByRoleId queries the menu tree of the specified role
func GetAllMenuByRoleId(c *gin.Context) {
	// bind current user role sorting. Hide specific users.
	user := GetCurrentUser(c)

	s := service.New(c)
	menus, ids, err := s.GetAllMenuByRoleId(&user.Role, utils.Str2Uint(c.Param("roleId")))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	var resp response.MenuTreeWithAccessResponseStruct
	resp.AccessIds = ids
	utils.Struct2StructByJson(menus, &resp.List)
	response.SuccessWithData(resp)
}

// GetMenus queries all menus.
func GetMenus(c *gin.Context) {
	// bind current user role sorting. Hide specific users.
	user := GetCurrentUser(c)

	s := service.New(c)
	menus := s.GetMenus(&user.Role)
	// convert to MenuTreeResponseStruct.
	var resp []response.MenuTreeResponseStruct
	utils.Struct2StructByJson(menus, &resp)
	response.SuccessWithData(resp)
}

func CreateMenu(c *gin.Context) {
	user := GetCurrentUser(c)
	// bind request body to struct.
	var req request.CreateMenuRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	// parameter validate.
	err = global.NewValidatorError(global.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// record current creator information
	req.Creator = user.Username

	s := service.New(c)
	err = s.CreateMenu(&user.Role, &req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

func UpdateMenuById(c *gin.Context) {
	// bind request body to struct.
	var req request.UpdateMenuRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	// eet the menuId in the path.
	menuId := utils.Str2Uint(c.Param("menuId"))
	if menuId == 0 {
		response.FailWithMsg("incorrect menu id")
		return
	}

	s := service.New(c)
	// update data.
	err = s.UpdateById(menuId, req, new(models.SysMenu))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteMenuByIds Batch delete menu data by menuId.
func BatchDeleteMenuByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("params binding failed, please check the data type")
		return
	}

	s := service.New(c)
	// delete data.
	err = s.DeleteByIds(req.GetUintIds(), new(models.SysMenu))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
