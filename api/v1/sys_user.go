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

// GetUserInfo get current user information.
func GetUserInfo(c *gin.Context) {
	user := GetCurrentUser(c)

	var resp response.UserInfoResponseStruct
	utils.Struct2StructByJson(user, &resp)
	resp.Roles = []string{
		user.Role.Keyword,
	}

	resp.NickName = user.Role.Name
	resp.RoleSort = *user.Role.Sort
	response.SuccessWithData(resp)
}

// GetUsers used to get a list of users.
func GetUsers(c *gin.Context) {
	// binding parameters.
	var req request.UserListRequestStruct
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	// bind current user role sorting. Hide specific users.
	user := GetCurrentUser(c)
	req.CurrentRole = user.Role

	s := service.New(c)
	users, err := s.GetUsers(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// convert to ResponseStruct, hide some fields.
	var respStruct []response.UserListResponseStruct
	utils.Struct2StructByJson(users, &respStruct)

	var resp response.PageData
	resp.PageInfo = req.PageInfo
	resp.List = respStruct
	response.SuccessWithData(resp)
}

// GetCurrentUser get the current request user information.
func GetCurrentUser(c *gin.Context) models.SysUser {
	user, exists := c.Get("user")
	var newUser models.SysUser
	if !exists {
		return newUser
	}
	u, _ := user.(models.SysUser)

	s := service.New(c)
	newUser, _ = s.GetUserById(u.Id)
	return newUser
}

// UpdateUserById update user information.
func UpdateUserById(c *gin.Context) {
	// bind request body to struct.
	var req request.UpdateUserRequestStruct
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	userId := utils.Str2Uint(c.Param("userId"))
	if userId == 0 {
		response.FailWithMsg("the userId is incorrect")
		return
	}

	user := GetCurrentUser(c)
	if userId == user.Id {
		if req.Status != nil && *req.Status == models.SysUserStatusDisabled {
			response.FailWithMsg("cannot disable itself")
			return
		}
		if req.RoleId != nil && user.RoleId != *req.RoleId {
			if *user.Role.Sort != models.SysRoleSuperAdminSort {
				response.FailWithMsg("unable to change your own role, if you need to change, please contact the superior role")
			} else {
				response.FailWithMsg("can't change Super Admin role")
			}
			return
		}
	}
	// By the way, modify the user's nickName.
	var role models.SysRole
	err = global.Mysql.Model(new(models.SysRole)).Where("id = ?", req.RoleId).First(&role).Error
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	req.NickName = role.Name

	s := service.New(c)
	// update data.
	err = s.UpdateById(userId, req, new(models.SysUser))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// CreateUser create a user.
func CreateUser(c *gin.Context) {
	user := GetCurrentUser(c)
	// bind request body to struct.
	var req request.CreateUserRequestStruct
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
	req.Creator = user.Username

	s := service.New(c)
	// convert password to ciphertext.
	req.Password = utils.GenPwd(req.Password)
	err = s.CreateUser(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteUserByIds used to delete users in bulk.
func BatchDeleteUserByIds(c *gin.Context) {
	var req request.Req
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMsg("parameter binding failed, please check the data type")
		return
	}

	user := GetCurrentUser(c)
	if utils.ContainsUint(req.GetUintIds(), user.Id) {
		response.FailWithMsg("can't delete itself")
		return
	}

	s := service.New(c)
	err = s.DeleteByIds(req.GetUintIds(), new(models.SysUser))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
