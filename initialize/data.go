package initialize

import (
	"errors"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"
	"strings"

	"gorm.io/gorm"
)

var (
	creator        = "系统自动创建"
	status         = uint(1)
	invisible      = uint(0)
	noBreadcrumb   = uint(0)
	autoDeploy     = uint(1)
	sysRoleDevSort = uint(20) //nolint:gomnd    // 一般开发人员排序
)

// Data 初始化数据
//
//nolint:funlen
//nolint:gocyclo
func Data() { //nolint:gocyclo
	if !global.Conf.System.InitData {
		return
	}

	// 1. 初始化角色
	newRoles := make([]*models.SysRole, 0)
	roles := []*models.SysRole{
		{
			Name:    "超级管理员",
			Keyword: "super",
			Desc:    "超级管理员",
		},
		{
			Name:    "研发人员",
			Keyword: "dev",
			Desc:    "研发人员",
			Sort:    &sysRoleDevSort,
		},
	}
	for i, role := range roles { //nolint:gocritic
		sort := uint(i)
		id := uint(i + 1)
		roles[i].Id = id
		oldRole := models.SysRole{}
		err := global.Mysql.Where("id = ?", id).First(&oldRole).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			role.Id = id
			role.Creator = creator
			role.Status = &status
			if role.Sort == nil {
				role.Sort = &sort
			}
			newRoles = append(newRoles, role)
		}
	}
	if len(newRoles) > 0 {
		global.Mysql.Create(&newRoles)
	}

	// 2. 初始化菜单
	menus := []*models.SysMenu{
		{
			Name:      "Dashboard", // 对于想让子菜单显示在上层不显示的父级菜单不设置名字
			Title:     "首页",
			Icon:      "home-svg|svg",
			Component: "/dashboard/analysis/index",
			Path:      "/dashboard",
			Roles:     roles,
		},
		{
			Name:      "Node",
			Title:     "终端管理",
			Icon:      "server-svg|svg",
			Path:      "/node",
			Redirect:  "/node/list",
			Component: "LAYOUT",
			Children: []*models.SysMenu{
				{
					Name:      "List",
					Icon:      "servers-svg|svg",
					Title:     "终端列表",
					Path:      "list",
					Component: "/node/list/index",
					Roles:     roles,
				},
				{
					Name:      "Ssh",
					Icon:      "terminal|svg",
					Title:     "远程连接",
					Path:      "ssh",
					Component: "/node/ssh/index",
				},
				{
					Name:      "My",
					Icon:      "like|svg",
					Title:     "我的收藏",
					Path:      "my",
					Component: "/node/my/index",
					Roles:     roles,
				},
			},
			Roles: roles,
		},
		{
			Name:      "Back",
			Title:     "后台管理",
			Icon:      "back-config-svg|svg",
			Path:      "/back",
			Component: "LAYOUT",
			Redirect:  "/back/label",
			Roles:     roles,
			Children: []*models.SysMenu{
				{
					Name:      "Label",
					Title:     "标签管理",
					Icon:      "labels-svg|svg",
					Path:      "label",
					Component: "/back/label/index",
					Roles:     roles,
				},
				{
					Name:      "Worker",
					Title:     "worker配置",
					Icon:      "worker-svg|svg",
					Path:      "worker",
					Component: "/back/worker/index",
					Roles:     roles,
				},
				{
					Name:      "Cron",
					Title:     "定时开关机",
					Icon:      "cron-svg|svg",
					Path:      "cron",
					Component: "/back/cron/index",
					Roles:     roles,
				},
			},
		},
		{
			Name:      "System",
			Title:     "系统管理",
			Icon:      "system-config-svg|svg",
			Path:      "/sys",
			Component: "LAYOUT",
			Redirect:  "/sys/role",
			Roles:     roles,
			Children: []*models.SysMenu{
				{
					Name:      "Role",
					Title:     "角色管理",
					Icon:      "role-svg|svg",
					Path:      "role",
					Component: "/sys/role/index",
					Roles:     roles,
				},
				{
					Name:      "User",
					Title:     "用户管理",
					Icon:      "users-svg|svg",
					Path:      "user",
					Component: "/sys/user/index",
					Roles:     roles,
				},
				{
					Name:      "Apis",
					Title:     "接口管理",
					Icon:      "api-svg|svg",
					Path:      "apis",
					Component: "/sys/apis/index",
				},
				{
					Name:      "Menu",
					Title:     "菜单管理",
					Icon:      "menus-svg|svg",
					Path:      "menu",
					Component: "/sys/menu/index",
				},
				{
					Name:      "OperationLog",
					Title:     "操作日志",
					Icon:      "operation-log|svg",
					Path:      "operation-log",
					Component: "/sys/operation-log/index",
				},
				{
					Name:      "About",
					Title:     "关于",
					Path:      "about",
					Visible:   &invisible,
					Component: "/sys/about/index",
					Roles:     roles,
				},
			},
		},
	}
	menus = genMenu(0, menus, roles[0])
	createMenu(menus)

	// 3. 初始化用户
	users := []models.SysUser{
		{
			Username:     "admin",
			Password:     utils.GenPwd("123456"),
			Email:        "admin@example.com",
			OfficeName:   "开发科",
			OrgName:      "开发部",
			WorkPlace:    "地球",
			NickName:     "超级管理员",
			Introduction: "超管就是不讲武德",
		},
		{
			Username:     "guest",
			Password:     utils.GenPwd("123456"),
			Email:        "guest@example.com",
			OfficeName:   "开发科",
			OrgName:      "开发部",
			WorkPlace:    "地球",
			NickName:     "开发人员",
			Introduction: "这个人很懒，什么也没留下",
		},
	}
	newUsers := make([]models.SysUser, 0)
	for i, user := range users { //nolint:gocritic
		id := uint(i + 1)
		oldUser := models.SysUser{}
		err := global.Mysql.Where("id = ?", id).First(&oldUser).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user.Id = id
			user.Creator = creator
			if user.RoleId == 0 {
				user.RoleId = id
			}
			newUsers = append(newUsers, user)
		}
	}
	if len(newUsers) > 0 {
		global.Mysql.Create(&newUsers)
	}

	// 4. 初始化接口
	apis := []models.SysApi{
		{
			Method:   "GET",
			Path:     "/v1/dashboard/countData",
			Category: "dashboard",
			Desc:     "获取首页统计数据",
		},
		{
			Method:   "GET",
			Path:     "/v1/dashboard/regionNodeData",
			Category: "dashboard",
			Desc:     "获取首页以地域分类的机器统计数据",
		},
		{
			Method:   "GET",
			Path:     "/v1/dashboard/managerNodeData",
			Category: "dashboard",
			Desc:     "获取首页以责任人分类的机器统计数据",
		},
		{
			Method:   "GET",
			Path:     "/v1/dashboard/performanceNodeData",
			Category: "dashboard",
			Desc:     "获取首页以性能分类的机器统计数据",
		},
		{
			Method:   "GET",
			Path:     "/v1/dashboard/healthNodeData",
			Category: "dashboard",
			Desc:     "获取首页以健康度分类的机器统计数据",
		},
		{
			Method:   "POST",
			Path:     "/v1/base/login",
			Category: "base",
			Desc:     "用户认证",
		},
		{
			Method:   "POST",
			Path:     "/v1/base/logout",
			Category: "base",
			Desc:     "用户认证",
		},
		{
			Method:   "POST",
			Path:     "/v1/base/refreshToken",
			Category: "base",
			Desc:     "用户token刷新",
		},
		{
			Method:   "GET",
			Path:     "/v1/base/idempotenceToken",
			Category: "base",
			Desc:     "获取幂等性token",
		},
		{
			Method:   "GET",
			Path:     "/v1/user/info",
			Category: "user",
			Desc:     "获取当前登录用户信息",
		},
		{
			Method:   "GET",
			Path:     "/v1/user/list",
			Category: "user",
			Desc:     "获取用户列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/user/create",
			Category: "user",
			Desc:     "创建用户",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/user/update/:userId",
			Category: "user",
			Desc:     "更新用户",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/user/delete/batch",
			Category: "user",
			Desc:     "删除用户",
		},
		{
			Method:   "GET",
			Path:     "/v1/menu/tree",
			Category: "menu",
			Desc:     "获取权限菜单",
		},
		{
			Method:   "GET",
			Path:     "/v1/menu/list",
			Category: "menu",
			Desc:     "获取菜单列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/menu/create",
			Category: "menu",
			Desc:     "创建菜单",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/menu/update/:menuId",
			Category: "menu",
			Desc:     "更新菜单",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/menu/delete/batch",
			Category: "menu",
			Desc:     "批量删除菜单",
		},
		{
			Method:   "GET",
			Path:     "/v1/role/list",
			Category: "role",
			Desc:     "获取角色列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/role/create",
			Category: "role",
			Desc:     "创建角色",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/role/update/:roleId",
			Category: "role",
			Desc:     "更新角色",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/role/delete/batch",
			Category: "role",
			Desc:     "批量删除角色",
		},
		{
			Method:   "GET",
			Path:     "/v1/label/list",
			Category: "label",
			Desc:     "获取标签列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/label/create",
			Category: "label",
			Desc:     "创建标签",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/label/update/:labelId",
			Category: "label",
			Desc:     "更新标签",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/label/delete/batch",
			Category: "label",
			Desc:     "批量删除标签",
		},
		{
			Method:   "GET",
			Path:     "/v1/menu/all/:roleId",
			Category: "menu",
			Desc:     "查询指定角色的菜单树",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/role/menus/update/:roleId",
			Category: "role",
			Desc:     "更新角色的权限菜单",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/role/apis/update/:roleId",
			Category: "role",
			Desc:     "更新角色的权限接口",
		},
		{
			Method:   "GET",
			Path:     "/v1/collect/my",
			Category: "collection",
			Desc:     "获取我收藏的机器信息",
		},
		{
			Method:   "POST",
			Path:     "/v1/collect/my",
			Category: "collection",
			Desc:     "添加机器到我的收藏",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/collect/delete/my",
			Category: "collection",
			Desc:     "取消我收藏的某个机器",
		},
		{
			Method:   "GET",
			Path:     "/v1/node/list",
			Category: "node",
			Desc:     "获取机器节点列表",
		},
		{
			Method:   "GET",
			Path:     "/v1/node/shell/ws",
			Category: "node",
			Desc:     "机器终端shell长连接",
		},
		{
			Method:   "GET",
			Path:     "/v1/node/vnc/ws",
			Category: "node",
			Desc:     "机器终端vnc长连接",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/node/shell/ws/resize",
			Category: "node",
			Desc:     "机器终端shell窗口变更",
		},
		{
			Method:   "POST",
			Path:     "/v1/node/shell/connect",
			Category: "node",
			Desc:     "测试机器终端shell连接",
		},
		{
			Method:   "GET",
			Path:     "/v1/node/shell/dir",
			Category: "node",
			Desc:     "获取机器终端shell路径下文件信息",
		},
		{
			Method:   "GET",
			Path:     "/v1/node/shell/file",
			Category: "node",
			Desc:     "获取机器终端shell路径下文件内容",
		},
		{
			Method:   "POST",
			Path:     "/v1/node/shell/file/download",
			Category: "node",
			Desc:     "下载机器终端shell路径下文件",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/node/shell/file/update",
			Category: "node",
			Desc:     "更新机器终端shell路径下文件",
		},
		{
			Method:   "POST",
			Path:     "/v1/node/create",
			Category: "node",
			Desc:     "创建机器节点",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/node/update/:nodeId",
			Category: "node",
			Desc:     "更新机器节点",
		},
		{
			Method:   "POST",
			Path:     "/v1/node/reboot",
			Category: "node",
			Desc:     "批量重启节点",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/node/refresh/:nodeId",
			Category: "node",
			Desc:     "刷新机器节点",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/node/delete/batch",
			Category: "node",
			Desc:     "批量删除机器节点",
		},
		{
			Method:   "GET",
			Path:     "/v1/api/list",
			Category: "api",
			Desc:     "获取接口列表",
		},
		{
			Method:   "GET",
			Path:     "/v1/api/all/category/:roleId",
			Category: "api",
			Desc:     "获取角色对应的接口列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/api/create",
			Category: "api",
			Desc:     "创建接口",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/api/update/:roleId",
			Category: "api",
			Desc:     "更新接口",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/api/delete/batch",
			Category: "api",
			Desc:     "批量删除接口",
		},
		{
			Method:   "GET",
			Path:     "/v1/upload/file",
			Category: "upload",
			Desc:     "获取文件块信息以及上传完成部分",
		},
		{
			Method:   "POST",
			Path:     "/v1/upload/file",
			Category: "upload",
			Desc:     "上传文件(分片)",
		},
		{
			Method:   "POST",
			Path:     "/v1/upload/merge",
			Category: "upload",
			Desc:     "合并分片文件",
		},
		{
			Method:   "GET",
			Path:     "/v1/operation/log/list",
			Category: "operation-log",
			Desc:     "获取操作日志列表",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/operation/log/delete/batch",
			Category: "operation-log",
			Desc:     "批量删除操作日志",
		},
		{
			Method:   "GET",
			Path:     "/v1/worker/list",
			Category: "worker",
			Desc:     "获取worker列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/worker/create",
			Category: "worker",
			Desc:     "创建worker",
		},
		{
			Method:   "PATCH",
			Path:     "/v1/worker/update/:workerId",
			Category: "worker",
			Desc:     "更新worker",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/worker/delete/batch",
			Category: "worker",
			Desc:     "批量删除worker",
		},
		{
			Method:   "GET",
			Path:     "/v1/cron/list",
			Category: "cron",
			Desc:     "获取定时开关机任务列表",
		},
		{
			Method:   "POST",
			Path:     "/v1/cron/create",
			Category: "cron",
			Desc:     "创建定时开关机任务",
		},
		{
			Method:   "POST",
			Path:     "/v1/cron/update/:userId",
			Category: "cron",
			Desc:     "更新定时开关机任务",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/cron/delete/batch",
			Category: "cron",
			Desc:     "删除定时开关机任务",
		},
		{
			Method:   "POST",
			Path:     "/v1/secure/container-report/:nodeId",
			Category: "secure",
			Desc:     "获取并更新机器节点容器的vul或sbom安全报告",
		},
		{
			Method:   "POST",
			Path:     "/v1/secure/bare-report/:nodeId",
			Category: "secure",
			Desc:     "获取并更新机器节点裸金属的vul或sbom安全报告",
		},
		{
			Method:   "GET",
			Path:     "/v1/secure/stats/:nodeId",
			Category: "secure",
			Desc:     "获取机器节点所有镜像",
		},
		{
			Method:   "GET",
			Path:     "/v1/secure/score/:nodeId",
			Category: "secure",
			Desc:     "获取机器节点的安全分数",
		},
		{
			Method:   "POST",
			Path:     "/v1/secure/bare/:nodeId",
			Category: "secure",
			Desc:     "修复裸金属的安全问题",
		},
		{
			Method:   "POST",
			Path:     "/v1/secure/container/:nodeId",
			Category: "secure",
			Desc:     "修复容器的安全问题",
		},
		{
			Method:   "POST",
			Path:     "/v1/tune/cleanup/:nodeId",
			Category: "tune",
			Desc:     "系统清理",
		},
		{
			Method:   "POST",
			Path:     "/v1/tune/auto/rollback/:nodeId",
			Category: "tune",
			Desc:     "系统调优还原",
		},
		{
			Method:   "POST",
			Path:     "/v1/tune/auto/set/:nodeId",
			Category: "tune",
			Desc:     "智能调优",
		},
		{
			Method:   "POST",
			Path:     "/v1/tune/scene/:nodeId",
			Category: "tune",
			Desc:     "设置场景调优",
		},
		{
			Method:   "POST",
			Path:     "/v1/tune/turbo/:nodeId",
			Category: "tune",
			Desc:     "系统性能加速",
		},
		{
			Method:   "GET",
			Path:     "/v1/tune/auto/list/:nodeId",
			Category: "tune",
			Desc:     "获取调优记录",
		},
		{
			Method:   "GET",
			Path:     "/v1/tune/score/:nodeId",
			Category: "tune",
			Desc:     "获取调优分数",
		},
		{
			Method:   "GET",
			Path:     "/v1/secure/risk/:nodeId",
			Category: "secure",
			Desc:     "获取机器节点风险条数",
		},
		{
			Method:   "DELETE",
			Path:     "/v1/tune/auto/delete",
			Category: "tune",
			Desc:     "删除智能调优记录",
		},
		{
			Method:   "POST",
			Path:     "/v1/secure/fix/:nodeId",
			Category: "secure",
			Desc:     "修复机器节点安全风险",
		},
	}
	newApis := make([]models.SysApi, 0)
	newRoleCasbins := make([]models.SysRoleCasbin, 0)
	for i, api := range apis { //nolint:gocritic
		id := uint(i + 1)
		oldApi := models.SysApi{}
		err := global.Mysql.Where("id = ?", id).First(&oldApi).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.Id = id
			api.Creator = creator
			newApis = append(newApis, api)
			// 超级管理员拥有所有API权限
			newRoleCasbins = append(newRoleCasbins, models.SysRoleCasbin{
				Keyword: roles[0].Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
			basePaths := []string{
				"/base/login",
				"/base/logout",
				"/base/refreshToken",
				"/base/idempotenceToken",
				"/user/info",
				"/menu/tree",
				"/dashboard/countData",
				"/dashboard/regionNodeData",
				"/dashboard/managerNodeData",
				"/dashboard/performanceNodeData",
				"/dashboard/healthNodeData",
				"/node/list",
			}
			p := strings.TrimPrefix(api.Path, "/"+global.Conf.System.ApiVersion)
			if utils.Contains(basePaths, p) {
				// 非超级管理员有基础权限
				for i := 1; i < len(roles); i++ {
					newRoleCasbins = append(newRoleCasbins, models.SysRoleCasbin{
						Keyword: roles[i].Keyword,
						Path:    api.Path,
						Method:  api.Method,
					})
				}
			}
		}
	}
	if len(newApis) > 0 {
		global.Mysql.Create(&newApis)
	}
	if len(newRoleCasbins) > 0 {
		s := service.New(nil)
		_, err := s.CreateRoleCasbins(newRoleCasbins)
		if err != nil {
			return
		}
	}

	// 5. 初始化worker
	workers := []*models.SysWorker{
		{
			Name:       "metalmetrics",
			Desc:       "用于获取机器节点的metrics信息",
			Port:       19091,
			AutoDeploy: &autoDeploy,
			DeployCmd: "curl -k -uadmin:123456 -L" +
				" https://example/factory/devops-metalflow/metalmetrics/metalmetrics.sh" +
				" -o metalmetrics.sh; chmod +x metalmetrics.sh",
			ReloadCmd:  `./metalmetrics.sh reload`,
			StopCmd:    `./metalmetrics.sh stop`,
			StartCmd:   `./metalmetrics.sh start`,
			ServiceReq: "metalmetrics/metrics",
			CheckReq:   "metalmetrics/version",
		},
		{
			Name:       "metaltask",
			Desc:       "用于接收远程文件并执行",
			Port:       19092,
			AutoDeploy: &autoDeploy,
			DeployCmd: "curl -k -uadmin:123456 -L" +
				" https://example/factory/devops-metalflow/metalmetrics/metalmetrics.sh" +
				" -o metaltask.sh; chmod +x metaltask.sh",
			ReloadCmd: `./metaltask.sh reload`,
			StopCmd:   `./metaltask.sh stop`,
			StartCmd:  `./metaltask.sh start`,
		},
		{
			Name:       "metalsecure",
			Desc:       "用于获取机器节点的安全报告",
			Port:       19094,
			AutoDeploy: &autoDeploy,
			DeployCmd: "curl -k -uadmin:123456 -L -L" +
				" https://example/factory/devops-metalflow/metalmetrics/metalmetrics.sh" +
				" -o metalsecure.sh; chmod +x metalsecure.sh",
			ReloadCmd: `./metalsecure.sh reload`,
			StopCmd:   `./metalsecure.sh stop`,
			StartCmd:  `./metalsecure.sh start`,
		},
		{
			Name:       "metaltune",
			Desc:       "用于对机器节点进行调优",
			Port:       19093,
			AutoDeploy: &autoDeploy,
			DeployCmd: "curl -k -uadmin:123456 -L" +
				" https://example/factory/devops-metalflow/metalmetrics/metalmetrics.sh" +
				" -o metaltune.sh; chmod +x metaltune.sh",
			ReloadCmd: `./metaltune.sh reload`,
			StopCmd:   `./metaltune.sh stop`,
			StartCmd:  `./metaltune.sh start`,
		},
	}
	newWorkers := make([]*models.SysWorker, 0)
	for i, worker := range workers {
		id := uint(i + 1)
		oldWorker := models.SysWorker{}
		err := global.Mysql.Where("id = ?", id).First(&oldWorker).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			worker.Id = id
			worker.Creator = creator
			newWorkers = append(newWorkers, worker)
		}
	}
	if len(newWorkers) > 0 {
		global.Mysql.Create(newWorkers)
	}

	// 6. 初始化系统调优的各场景profile
	tuneScenes := []*models.SysNodeTuneScene{
		{
			Scene:   "build",
			Profile: "hello world",
		},
		{
			Scene:   "program",
			Profile: "hello world",
		},
		{
			Scene:   "ci",
			Profile: "hello world",
		},
		{
			Scene:   "k8s",
			Profile: "hello world",
		},
		{
			Scene:   "ai",
			Profile: "hello world",
		},
		{
			Scene:   "test",
			Profile: "hello world",
		},
		{
			Scene:   "origin",
			Profile: "hello world",
		},
	}
	newTuneScenes := make([]*models.SysNodeTuneScene, 0)
	for i, tuneScene := range tuneScenes {
		id := uint(i + 1)
		oldTuneScene := models.SysNodeTuneScene{}
		err := global.Mysql.Where("id = ?", id).First(&oldTuneScene).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tuneScene.Id = id
			tuneScene.Creator = creator
			newTuneScenes = append(newTuneScenes, tuneScene)
		}
	}
	if len(newTuneScenes) > 0 {
		global.Mysql.Create(newTuneScenes)
	}
}

var menuTotal = 0

// 生成菜单
func genMenu(parentId uint, menus []*models.SysMenu, superRole *models.SysRole) []*models.SysMenu {
	newMenus := make([]*models.SysMenu, len(menus))
	// sort
	for i, menu := range menus { //nolint:gocritic
		sort := uint(i)
		menu.Sort = &sort
		menu.Creator = creator
		menu.Status = &status
		newMenus[i] = menu
	}
	// id
	for i, menu := range newMenus { //nolint:gocritic
		menuTotal++
		menu.Id = uint(menuTotal)
		newMenus[i] = menu
	}
	// children
	for i, menu := range newMenus { //nolint:gocritic
		menu.Children = genMenu(menu.Id, menu.Children, superRole)
		newMenus[i] = menu
	}
	// parentId
	for i, menu := range newMenus { //nolint:gocritic
		if parentId > 0 {
			menu.ParentId = parentId
		} else {
			menu.Breadcrumb = &noBreadcrumb // 面包屑不可见
		}
		if menu.Roles == nil {
			menu.Roles = []*models.SysRole{
				superRole,
			}
		}
		newMenus[i] = menu
	}
	return newMenus
}

// 将生成的菜单保存到数据库
func createMenu(menus []*models.SysMenu) {
	for _, menu := range menus {
		oldMenu := models.SysMenu{}
		err := global.Mysql.Where("id = ?", menu.Id).First(&oldMenu).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Mysql.Create(&menu) //nolint:gosec
		}
		if len(menu.Children) > 0 {
			createMenu(menu.Children)
		}
	}
}
