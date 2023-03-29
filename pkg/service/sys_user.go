package service

import (
	"errors"
	"fmt"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/utils"
	"strings"
)

// LoginCheck checks user login information verification.
func (s *MysqlService) LoginCheck(req *request.UserAuthRequestStruct) (*models.SysUser, error) {
	var user models.SysUser
	err := s.TX.Preload("Role").Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}
	if ok := utils.ComparePwd(req.Password, user.Password); !ok {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}
	return &user, nil
}

// GetUsers 获取用户
func (s *MysqlService) GetUsers(req *request.UserListRequestStruct) ([]models.SysUser, error) {
	var err error
	list := make([]models.SysUser, 0)
	db := global.Mysql.
		Model(&models.SysUser{}).
		Order("created_at DESC")
	// 非超级管理员
	if *req.CurrentRole.Sort != models.SysRoleSuperAdminSort {
		roleIds, err := s.GetRoleIdsBySort(*req.CurrentRole.Sort) //nolint:govet
		if err != nil {
			return list, err
		}
		db = db.Where("role_id IN (?)", roleIds)
	}
	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	officeName := strings.TrimSpace(req.OfficeName)
	if officeName != "" {
		db = db.Where("office_name LIKE ?", fmt.Sprintf("%%%s%%", officeName))
	}
	orgName := strings.TrimSpace(req.OrgName)
	if orgName != "" {
		db = db.Where("org_name LIKE ?", fmt.Sprintf("%%%s%%", orgName))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	if req.Status != nil {
		if *req.Status > 0 {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
	}
	// 查询条数
	err = s.Find(db, &req.PageInfo, &list)
	return list, err
}

// CreateUser 创建用户
func (s *MysqlService) CreateUser(req *request.CreateUserRequestStruct) (err error) {
	user := new(models.SysUser)
	err = s.Create(req, &user)
	return
}

// GetUserById 获取单个用户
func (s *MysqlService) GetUserById(id uint) (models.SysUser, error) {
	var user models.SysUser

	err := s.TX.Preload("Role").
		Where("id = ?", id).
		// 状态为正常
		Where("status = ?", models.SysUserStatusNormal).
		First(&user).Error
	return user, err
}
