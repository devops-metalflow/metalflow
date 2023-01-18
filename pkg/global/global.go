package global

import (
	"embed"
	"errors"
	"metalflow/pkg/async"
	"metalflow/pkg/cron"
	"os"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

var (
	// Conf 系统配置
	Conf Configuration
	// ConfBox packr盒子用于打包配置文件到golang编译后的二进制程序中
	ConfBox *CustomConfBox
	// Log zap日志
	Log *zap.SugaredLogger
	// Mysql Postgresql实例
	Mysql *gorm.DB
	// Redis redis实例
	Redis *redis.Client
	// Machinery machinery实例
	Machinery *async.Machinery
	// CasbinEnforcer cabin实例
	CasbinEnforcer *casbin.Enforcer
	// Validate validation.v9校验器
	Validate *validator.Validate
	// Translator validation.v9相关翻译器
	Translator ut.Translator
	// 定时任务管理器
	Cron *cron.Client
)

// CustomConfBox 自定义配置盒子
type CustomConfBox struct {
	// 配置文件路径环境变量
	ConfFile string
	// 丢弃packr盒子，docker环境下部署有报错 使用go1.16新特性embed来嵌入配置文件到二进制执行文件中
	EmbedFs *embed.FS
	// viper实例
	ViperIns *viper.Viper
}

// Find 查找指定配置
func (c *CustomConfBox) Find(filename string) []byte {
	bs, _ := os.ReadFile(filename)
	if len(bs) == 0 {
		bs, _ = c.EmbedFs.ReadFile(filename)
	}
	return bs
}

// NewValidatorError 只返回一个错误即可
func NewValidatorError(err error, custom map[string]string) (e error) {
	if err == nil {
		return
	}
	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		tranStr := e.Translate(Translator)
		// 判断错误字段是否在自定义集合中，如果在，则替换错误信息中的字段
		if v, ok := custom[e.Field()]; ok {
			return errors.New(strings.Replace(tranStr, e.Field(), v, 1))
		} else {
			return errors.New(tranStr)
		}
	}
	return
}

// GetTx 获取事务对象
func GetTx(c *gin.Context) *gorm.DB {
	// 默认使用无事务的postgresql
	tx := Mysql
	if c != nil {
		method := ""
		if c.Request != nil {
			method = c.Request.Method
		}
		if !(method == "OPTIONS" || method == "GET" || !Conf.System.Transaction) {
			// 从context对象中读取事务对象
			txKey, exists := c.Get("tx")
			if exists {
				if item, ok := txKey.(*gorm.DB); ok {
					tx = item
				}
			}
		}
	}
	return tx
}
