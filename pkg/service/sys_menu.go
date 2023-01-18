package service

import (
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/utils"

	"github.com/thedevsaddam/gojsonq/v2"
)

// GetMenuTree 获取权限菜单树
func (s *MysqlService) GetMenuTree(roleId uint) ([]*models.SysMenu, error) {
	var role models.SysRole
	err := s.TX.Table(new(models.SysRole).TableName()).Preload("Menus").Where("id = ?", roleId).Find(&role).Error
	menus := make([]*models.SysMenu, 0)
	if err != nil {
		return menus, err
	}
	// 生成菜单树
	tree := GenMenuTree(nil, role.Menus)
	return tree, nil
}

// GetMenus 获取所有菜单
func (s *MysqlService) GetMenus(currentRole *models.SysRole) []*models.SysMenu {
	menus := s.getAllMenu(currentRole)
	// 生成菜单树
	tree := GenMenuTree(nil, menus)
	return tree
}

// GenMenuTree 生成菜单树, 主要是让menu的children归位
func GenMenuTree(parent *models.SysMenu, menus []*models.SysMenu) []*models.SysMenu {
	tree := make([]*models.SysMenu, 0)
	// parentId默认为0, 表示根菜单
	var parentId uint
	if parent != nil {
		parentId = parent.Id
	} else {
		// 将菜单转为json再排序
		newMenus := make([]*models.SysMenu, 0)
		list := gojsonq.New().FromString(utils.Struct2Json(menus)).SortBy("sort").Get()
		// 再转为json
		utils.Struct2StructByJson(list, &newMenus)
		menus = newMenus
	}

	for _, menu := range menus {
		// 父菜单编号一致
		if menu.ParentId == parentId {
			// 递归获取子菜单
			menu.Children = GenMenuTree(menu, menus)
			// 加入菜单树
			tree = append(tree, menu)
		}
	}
	return tree
}

// GetAllMenuByRoleId 根据权限编号获取全部菜单
func (s *MysqlService) GetAllMenuByRoleId(currentRole *models.SysRole, roleId uint) ([]*models.SysMenu, []uint, error) {
	// 有权限访问的id列表
	accessIds := make([]uint, 0)
	// 查询全部菜单
	allMenu := s.getAllMenu(currentRole)
	// 查询角色拥有菜单
	roleMenus := s.getRoleMenus(roleId)
	// 生成菜单树
	tree := GenMenuTree(nil, allMenu)
	// 获取id列表
	for _, menu := range roleMenus {
		accessIds = append(accessIds, menu.Id)
	}
	// 只保留选中项目
	accessIds = models.GetCheckedMenuIds(accessIds, allMenu)
	return tree, accessIds, nil
}

// CreateMenu 创建菜单
func (s *MysqlService) CreateMenu(currentRole *models.SysRole, req *request.CreateMenuRequestStruct) (err error) {
	var menu models.SysMenu
	utils.Struct2StructByJson(req, &menu)
	// 创建数据
	err = s.TX.Create(&menu).Error
	if err != nil {
		return err
	}
	// 自己创建的菜单需绑定权限
	menuReq := request.UpdateIncrementalIdsRequestStruct{
		Create: []uint{menu.Id},
	}
	err = s.UpdateRoleMenusById(currentRole, currentRole.Id, menuReq)
	return
}

// 获取权限菜单, 非菜单树
func (s *MysqlService) getRoleMenus(roleId uint) []*models.SysMenu {
	var role models.SysRole
	// 根据权限编号获取菜单
	err := s.TX.Preload("Menus").Where("id = ?", roleId).First(&role).Error
	global.Log.Warn("[getRoleMenu]", err)
	return role.Menus
}

// 获取全部菜单, 非菜单树
func (s *MysqlService) getAllMenu(currentRole *models.SysRole) []*models.SysMenu {
	menus := make([]*models.SysMenu, 0)
	// 查询关系表
	relations := make([]models.RelationRoleMenu, 0)
	menuIds := make([]uint, 0)
	query := s.TX.Model(models.RelationRoleMenu{})
	var err error
	// 非超级管理员
	if *currentRole.Sort != models.SysRoleSuperAdminSort {
		query = query.Where("sys_role_id = ?", currentRole.Id)
		err = query.Find(&relations).Error
		if err != nil {
			return menus
		}
		for _, relation := range relations {
			menuIds = append(menuIds, relation.SysMenuId)
		}
		// 查询所有菜单
		err = s.TX.Order("sort").Where("id IN (?)", menuIds).Find(&menus).Error
	} else {
		err = s.TX.Order("sort").Find(&menus).Error
	}

	global.Log.Warn("[getAllMenu]", err)
	return menus
}
