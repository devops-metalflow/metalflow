package tests

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/response"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func GetMock() sqlmock.Sqlmock {
	// 每一个测试用例都需要重新new sqlmock实例，不然会有测试混淆报错
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		fmt.Printf("init sqlmock failed: %v\n", err)
		panic(err)
	}
	// 结合gorm sqlmock(测试时不需要真正连接数据库)
	db, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      sqlDB,
	}), &gorm.Config{})
	if err != nil {
		fmt.Printf("init gorm with sqlmock failed:%v\n", err)
		panic(err)
	}
	db = db.Debug()
	global.Mysql = db
	// set table name prefix.
	global.Conf.Mysql.TablePrefix = "tb"
	return mock
}

func SetLog() {
	// set log
	logger, _ := zap.NewProduction()
	global.Log = logger.Sugar()
}

// SetConfig 初始化配置文件
func SetConfig(cfgPath string) {
	// 初始化配置盒子
	var box global.CustomConfBox

	// 获取viper实例
	box.ViperIns = viper.New()
	global.ConfBox = &box
	v := box.ViperIns

	// 读取开发环境配置作为默认配置项
	readConfig(v, cfgPath)
	// 转换为结构体
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v, 配置文件: %s", err, global.ConfBox.ConfFile))
	}

	if global.Conf.System.ConnectTimeout < 1 {
		global.Conf.System.ConnectTimeout = 10
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
	v.SetConfigType("yml")
	config := global.ConfBox.Find(configFile)
	if len(config) == 0 {
		panic(fmt.Sprintf("初始化配置文件失败: %v", configFile))
	}
	// 加载配置
	if err := v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}
}

func SetCasbinEnforcer(cfgPath string) {
	e, err := casbin.NewEnforcer(cfgPath, false)
	if err != nil {
		fmt.Printf("init test casbinEnforcer failed:%v\n", err)
		panic(err)
	}
	global.CasbinEnforcer = e
}

var (
	CurrentUserId = uint(2) // nolint:gomnd
)

func MockGetCurrentUser(mock sqlmock.Sqlmock, choice uint) {
	// currentUser
	switch choice {
	case uint(0): // nolint:gomnd
		mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").
			WithArgs(CurrentUserId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(CurrentUserId))
	case uint(1): // nolint:gomnd
		mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").
			WithArgs(CurrentUserId, 1).
			WillReturnError(gorm.ErrRecordNotFound)
	case uint(2): // nolint:gomnd
		mock.ExpectQuery("SELECT (.*) FROM `tb_sys_user`").
			WithArgs(CurrentUserId, 1).
			WillReturnError(errors.New("db error"))
	default:
	}
}

func SetValidate() {
	// 实例化需要转换的语言, 中文
	chinese := zh.New()
	uni := ut.New(chinese, chinese)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()

	// 注册转换的语言为默认语言
	_ = translations.RegisterDefaultTranslations(validate, trans)

	global.Validate = validate
	global.Translator = trans
}

func SetContextUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user", models.SysUser{
			Model:    models.Model{Id: CurrentUserId},
			Username: "12345678",
			RoleId:   1,
			Role: models.SysRole{
				Model:   models.Model{Id: 1},
				Keyword: "tester",
			},
		})
	}
}

func RemoveContextUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Keys = nil
	}
}

func exception(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[Exception]未知异常: %v\n堆栈信息: %v", err, string(debug.Stack()))
			// 服务器异常
			resp := response.Resp{
				Code:   response.InternalServerError,
				Result: map[string]any{},
				Msg:    response.CustomError[response.InternalServerError],
			}
			// 以json方式写入响应
			response.JSON(c, http.StatusOK, resp)
		}
	}()
	c.Next()
}

// transaction 全局事务处理中间件
func transaction(c *gin.Context) {
	method := c.Request.Method
	noTransaction := false
	if method == "OPTIONS" || method == "GET" || !global.Conf.System.Transaction {
		// Options/GET方法 以及 未配置事务时不创建事务
		noTransaction = true
	}
	defer func() {
		// 获取事务对象
		tx := global.GetTx(c)
		if err := recover(); err != nil {
			if resp, ok := err.(response.Resp); ok {
				if !noTransaction {
					if resp.Code == response.Ok {
						// 有效的请求，提交事务
						tx.Commit()
					} else {
						// 回滚事务
						tx.Rollback()
					}
				}
				// 以json方式写入响应
				response.JSON(c, http.StatusOK, resp)
				c.Abort()
				return
			}
			if !noTransaction {
				// 回滚事务
				tx.Rollback()
			}
			// 继续向上层抛出异常
			panic(err)
		} else if !noTransaction {
			// 没有异常, 提交事务
			tx.Commit()
		}
		// 结束请求, 避免二次调用
		c.Abort()
	}()
	if !noTransaction {
		// 开启事务, 写入当前请求
		tx := global.Mysql.Begin()
		c.Set("tx", tx)
	}
	// 处理请求
	c.Next()
}

func GetRouter() *gin.Engine {
	r := gin.Default()
	r.Use(exception)
	r.Use(transaction)
	r.Use(SetContextUser())
	return r
}
