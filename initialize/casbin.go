package initialize

import (
	"fmt"
	"metalflow/pkg/global"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

const casbinConfigPath = "conf/rbac_model.conf"

// CasbinEnforcer 初始化cabin
func CasbinEnforcer() {
	e, err := mysqlCasbin()
	if err != nil {
		panic(fmt.Sprintf("初始化casbin策略管理器: %v", err))
	}
	global.CasbinEnforcer = e
	global.Log.Info("初始化casbin策略管理器完成")
}

func mysqlCasbin() (*casbin.Enforcer, error) {
	// casbin默认表名casbin_rule, 为了与项目统一改写一下规则
	// 注意: gormadapter.CasbinTableName内部添加了下划线, 这里不再多此一举
	a, err := gormadapter.NewAdapterByDBUseTableName(global.Mysql, global.Conf.Mysql.TablePrefix, "sys_casbin")
	if err != nil {
		return nil, err
	}
	// 加锁避免并发多次初始化cabinModel. 读取配置文件
	config := global.ConfBox.Find(casbinConfigPath)
	cabinModel := model.NewModel()
	// 从字符串中加载casbin配置
	err = cabinModel.LoadModelFromText(string(config))
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(cabinModel, a)
	if err != nil {
		return nil, err
	}
	// 加载策略
	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return e, nil
}
