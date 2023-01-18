package service

import (
	"fmt"
	"gorm.io/gorm"
	"metalflow/models"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/utils"
	"strings"
)

// GetApis 获取所有接口
func (s *MysqlService) GetApis(req *request.ApiRequestStruct) ([]models.SysApi, error) {
	var err error
	list := make([]models.SysApi, 0)
	query := s.TX.
		Model(&models.SysApi{}).
		Order("created_at DESC")
	method := strings.TrimSpace(req.Method)
	if method != "" {
		query = query.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		query = query.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		query = query.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		query = query.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	// 查询列表
	err = s.Find(query, &req.PageInfo, &list)
	return list, err
}

// GetAllApiGroupByCategoryByRoleId 根据权限编号获取以api分类分组的权限接口
//
//nolint:gocyclo
func (s *MysqlService) GetAllApiGroupByCategoryByRoleId(currentRole *models.SysRole, roleId uint) ([]response.ApiGroupByCategoryResponseStruct, []uint, error) { //nolint:lll
	// 接口树
	tree := make([]response.ApiGroupByCategoryResponseStruct, 0)
	// 有权限访问的id列表
	accessIds := make([]uint, 0)
	allApi := make([]*models.SysApi, 0)
	// 查询全部api
	err := s.TX.Find(&allApi).Error
	if err != nil {
		return tree, accessIds, err
	}
	var currentRoleId uint
	// 非超级管理员
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		currentRoleId = currentRole.Id
	}
	// 查询当前角色拥有api访问权限的casbin规则
	currentCasbins, err := s.GetCasbinListByRoleId(currentRoleId)
	if err != nil {
		return tree, accessIds, err
	}
	// 查询指定角色拥有api访问权限的casbin规则.当前角色只能在自己权限范围内操作, 不得越权
	casbins, err := s.GetCasbinListByRoleId(roleId)
	if err != nil {
		return tree, accessIds, err
	}

	// 找到当前角色的全部api
	newApi := make([]*models.SysApi, 0)
	for _, api := range allApi {
		path := api.Path
		method := api.Method
		for _, currentCasbin := range currentCasbins {
			// 该api有权限
			if path == currentCasbin.V1 && method == currentCasbin.V2 {
				newApi = append(newApi, api)
				break
			}
		}
	}

	// 通过分类进行分组归纳
	for _, api := range newApi {
		category := api.Category
		path := api.Path
		method := api.Method
		access := false
		for _, casbin := range casbins {
			// 该api有权限
			if path == casbin.V1 && method == casbin.V2 {
				access = true
				break
			}
		}
		// 加入权限集合
		if access {
			accessIds = append(accessIds, api.Id)
		}
		// 生成接口树
		existIndex := -1
		children := make([]response.ApiListResponseStruct, 0)
		for index, leaf := range tree {
			if leaf.Category == category {
				children = leaf.Children
				existIndex = index
				break
			}
		}
		// api结构转换
		var item response.ApiListResponseStruct
		utils.Struct2StructByJson(api, &item)
		item.Title = fmt.Sprintf("%s %s[%s]", item.Desc, item.Path, item.Method)
		children = append(children, item)
		if existIndex != -1 {
			// 更新元素
			tree[existIndex].Children = children
		} else {
			// 新增元素
			tree = append(tree, response.ApiGroupByCategoryResponseStruct{
				Title:    category + "分组",
				Category: category,
				Children: children,
			})
		}
	}
	return tree, accessIds, err
}

// CreateApi 创建接口
func (s *MysqlService) CreateApi(req *request.CreateApiRequestStruct) (err error) {
	api := new(models.SysApi)
	err = s.Create(req, &api)
	if err != nil {
		return err
	}
	// 添加了角色
	if len(req.RoleIds) > 0 {
		// 查询角色关键字
		var roles []*models.SysRole
		err = s.TX.Where("id IN (?)", req.RoleIds).Find(&roles).Error
		if err != nil {
			return
		}
		// 构建casbin规则
		cs := make([]models.SysRoleCasbin, 0)
		for _, role := range roles {
			cs = append(cs, models.SysRoleCasbin{
				Keyword: role.Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
		}
		// 批量创建
		_, err = s.BatchCreateRoleCasbins(cs)
	}
	return
}

// UpdateApiById 更新接口
func (s *MysqlService) UpdateApiById(id uint, req request.UpdateApiRequestStruct) (err error) {
	var api models.SysApi
	query := s.TX.Model(&api).Where("id = ?", id).First(&api)
	if query.Error == gorm.ErrRecordNotFound {
		return fmt.Errorf("记录不存在")
	}

	// 比对增量字段
	m := make(map[string]any, 0)
	utils.CompareDifferenceStruct2SnakeKeyByJson(api, req, &m)

	// 记录update前的旧数据, 执行Updates后api会变成新数据
	oldApi := api
	// 更新指定列
	err = query.Updates(m).Error

	// 对比api发生了哪些变化
	diff := make(map[string]any, 0)
	utils.CompareDifferenceStruct2SnakeKeyByJson(oldApi, api, &diff)

	path, ok1 := diff["path"]
	method, ok2 := diff["method"]
	if (ok1 && path != "") || (ok2 && method != "") {
		// path或method变化, 需要更新casbin规则
		// 查找当前接口都有哪些角色在使用
		oldCasbins := s.GetRoleCasbins(models.SysRoleCasbin{
			Path:   oldApi.Path,
			Method: oldApi.Method,
		})
		if len(oldCasbins) > 0 {
			keywords := make([]string, 0)
			for _, oldCasbin := range oldCasbins {
				keywords = append(keywords, oldCasbin.Keyword)
			}
			// 删除旧规则, 添加新规则
			_, _ = s.BatchDeleteRoleCasbins(oldCasbins)
			// 构建新casbin规则
			newCasbins := make([]models.SysRoleCasbin, 0)
			for _, keyword := range keywords {
				newCasbins = append(newCasbins, models.SysRoleCasbin{
					Keyword: keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
			// 批量创建
			_, err = s.BatchCreateRoleCasbins(newCasbins)
		}
	}
	return err
}

// DeleteApiByIds 批量删除接口
func (s *MysqlService) DeleteApiByIds(ids []uint) (err error) {
	var list []*models.SysApi
	query := s.TX.Where("id IN (?)", ids).Find(&list)
	if query.Error != nil {
		return query.Error
	}
	// 查找当前接口都有哪些角色在使用
	casbins := make([]models.SysRoleCasbin, 0)
	for _, api := range list {
		casbins = append(casbins, s.GetRoleCasbins(models.SysRoleCasbin{
			Path:   api.Path,
			Method: api.Method,
		})...)
	}
	// 删除所有规则
	_, _ = s.BatchDeleteRoleCasbins(casbins)
	return query.Delete(&models.SysApi{}).Error
}
