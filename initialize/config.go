package initialize

import (
	"bytes"
	"embed"
	"fmt"
	"metalflow/pkg/global"
	"metalflow/pkg/utils"
	"os"
	"strings"

	"github.com/spf13/viper"
)

//go:embed conf/config.dev.yml
//go:embed conf/config.prod.yml
//go:embed conf/config.test.yml
//go:embed conf/rbac_model.conf
var f embed.FS

const (
	configType            = "yml"
	developmentConfig     = "conf/config.dev.yml"
	productionConfig      = "conf/config.prod.yml"
	testConfig            = "conf/config.test.yml"
	defaultConnectTimeout = 5
)

// Config 初始化配置文件
func Config(filename string) {
	// 初始化配置盒子
	var box global.CustomConfBox
	// 读取命令行配置文件并判断
	if filename != "" {
		if strings.HasPrefix(filename, "/") {
			// 指定的目录为绝对路径
			box.ConfFile = filename
		} else {
			// 指定的目录为相对路径
			box.ConfFile = utils.GetWorkDir() + "/" + filename
		}
	}

	// 获取viper实例
	box.ViperIns = viper.New()
	// 得到go embed实例
	box.EmbedFs = &f
	global.ConfBox = &box
	v := box.ViperIns

	// 读取开发环境配置作为默认配置项
	readConfig(v, developmentConfig)
	// 将default中的配置全部以默认配置写入
	settings := v.AllSettings()
	for index, setting := range settings {
		v.SetDefault(index, setting)
	}
	// 根据环境变量再读取一次配置覆盖,因前端使用的VITE打包，故环境变量必须得以VITE开头
	env := strings.ToLower(os.Getenv("VITE_WEB_MODE"))
	configName := ""
	if env == global.Prod {
		configName = productionConfig
	} else if env == global.Test {
		configName = testConfig
	}
	if configName != "" {
		readConfig(v, configName)
	}
	// 最后读取命令行配置文件内容
	if box.ConfFile != "" {
		// 读取不同配置文件中的差异部分
		readConfig(v, box.ConfFile)
	}
	// 转换为结构体
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v, 配置文件: %s", err, global.ConfBox.ConfFile))
	}

	if global.Conf.System.ConnectTimeout < 1 {
		global.Conf.System.ConnectTimeout = defaultConnectTimeout
	}

	if strings.TrimSpace(global.Conf.System.UrlPathPrefix) == "" {
		global.Conf.System.UrlPathPrefix = "api"
	}

	if strings.TrimSpace(global.Conf.System.ApiVersion) == "" {
		global.Conf.System.UrlPathPrefix = "v1"
	}

	// 表前缀去掉后缀_
	if strings.TrimSpace(global.Conf.Mysql.TablePrefix) != "" && strings.HasSuffix(global.Conf.Mysql.TablePrefix, "_") {
		global.Conf.Mysql.TablePrefix = strings.TrimSuffix(global.Conf.Mysql.TablePrefix, "_")
	}
}

func readConfig(v *viper.Viper, configFile string) {
	v.SetConfigType(configType)
	config := global.ConfBox.Find(configFile)
	if len(config) == 0 {
		panic(fmt.Sprintf("初始化配置文件失败: %v", configFile))
	}
	// 加载配置
	if err := v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}
}
