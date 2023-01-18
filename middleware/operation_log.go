package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/casbin/casbin/v2/util"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/datatypes"
	"io"
	v1 "metalflow/api/v1"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	UnLogin    = "未登录"
	NoResponse = "{}"
)

var (
	// 定期缓存, 避免每次频繁查询数据库
	// nolint:gomnd
	apiCache = cache.New(24*time.Hour, 48*time.Hour)
	// 日志缓存
	logCache = make([]models.SysOperationLog, 0)
	logLock  sync.Mutex
)

// OperationLog for record operation log from every request.
// nolint:funlen
func OperationLog(c *gin.Context) { // nolint:gocyclo
	// 开始时间
	startTime := time.Now()
	// 读取body参数
	var body []byte
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		global.Log.Error(c, "读取请求体失败: %v", err)
	} else {
		// gin参数只能读取一次, 这里将其回写, 否则c.Next中的接口无法读取
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	// 避免服务器出现异常, 这里用defer保证一定可以执行
	defer func() {
		// GET/OPTIONS请求比较频繁无需写入日志
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodOptions {
			return
		}
		// 用户自定义请求无需写入日志
		for _, s := range global.Conf.System.OperationLogDisabledPathArr {
			if strings.Contains(c.Request.URL.Path, s) {
				return
			}
		}

		// 结束时间
		endTime := time.Now()

		if len(body) == 0 {
			body = []byte("{}")
		}
		contentType := c.Request.Header.Get("Content-Type")
		// 二进制文件类型需要特殊处理
		if strings.Contains(contentType, "multipart/form-data") {
			contentTypeArr := strings.Split(contentType, "; ")
			if len(contentTypeArr) == 2 { // nolint:gomnd
				// 读取boundary
				boundary := strings.TrimPrefix(contentTypeArr[1], "boundary=")
				// 通过multipart读取body参数全部内容
				b := strings.NewReader(string(body))
				r := multipart.NewReader(b, boundary)
				f, _ := r.ReadForm(int64(global.Conf.Upload.SingleMaxSize) << 20) // nolint:gomnd
				defer func(f *multipart.Form) {
					_ = f.RemoveAll()
				}(f)
				// 获取全部参数值
				params := make(map[string]string, 0)
				for key, val := range f.Value {
					// 保留第一个值就行了
					if len(val) > 0 {
						params[key] = val[0]
					}
				}
				params["content-type"] = "multipart/form-data"
				params["file"] = "二进制数据被忽略"
				// 将其转为json
				body = []byte(utils.Struct2Json(params))
			}
		}
		// 记录header
		header := make(map[string]string, 0)
		for k, v := range c.Request.Header {
			header[k] = strings.Join(v, " | ")
		}
		headerJson, _ := json.Marshal(header)
		log := models.SysOperationLog{
			Model: models.Model{
				// 记录最后时间
				CreatedAt: models.LocalTime{
					Time: endTime,
				},
			},
			// 请求方式
			Method: c.Request.Method,
			// 请求路径 去除url前缀
			Path: strings.TrimPrefix(c.Request.URL.Path, "/"+global.Conf.System.UrlPathPrefix),
			// 请求头
			Header: headerJson,
			// 请求体
			Body: datatypes.JSON(body),
			// 请求耗时
			Latency: endTime.Sub(startTime).Milliseconds(),
			// 浏览器标识
			UserAgent: c.Request.UserAgent(),
		}

		// 清理事务
		c.Set("tx", "")
		// 获取当前登录用户
		user := v1.GetCurrentUser(c)

		// 用户名
		if user.Id > 0 {
			log.UserName = user.Username
			log.RoleName = user.Role.Name
		} else {
			log.UserName = UnLogin
			log.RoleName = UnLogin
		}

		log.ApiDesc = getApiDesc(c, log.Method, log.Path)

		// 响应状态码
		log.Status = c.Writer.Status()
		// 响应数据
		resp, exists := c.Get(global.Conf.System.OperationLogKey)
		var data datatypes.JSON
		if exists {
			data, _ = json.Marshal(resp)
			// 是自定义的响应类型
			if item, ok := resp.(response.Resp); ok {
				// 未登录操作信息无需写入日志
				if item.Code == response.Unauthorized {
					return
				}
				log.Status = item.Code
			}
		} else {
			data = datatypes.JSON(NoResponse)
		}
		// gzip压缩
		log.Data = data
		// 记录操作时间
		log.CreatedAt = models.LocalTime{
			Time: time.Now(),
		}
		// 操作日志晚点写入数据库，目前访问量不大，每10条存入到数据库一次
		logLock.Lock()
		logCache = append(logCache, log)
		if len(logCache) >= 10 { // nolint:gomnd
			logs := logCache
			go global.Mysql.Create(&logs)
			logCache = make([]models.SysOperationLog, 0)
		}
		logLock.Unlock()
	}()
	c.Next()
}

// getApiDesc return api desc.
func getApiDesc(c *gin.Context, method, path string) string {
	desc := "无"
	apiMap := make(map[string][]models.SysApi, 0)
	oldCache1, ok1 := apiCache.Get(fmt.Sprintf("%s_%s", method, path))
	if ok1 {
		desc, _ = oldCache1.(string)
		return desc
	}
	oldCache2, ok2 := apiCache.Get("apiMap")
	if ok2 {
		apiMap, _ = oldCache2.(map[string][]models.SysApi)
	} else {
		// 获取当前接口
		s := service.New(c)
		apis, err := s.GetApis(&request.ApiRequestStruct{
			PageInfo: response.PageInfo{
				NoPagination: true,
			},
		})
		if err == nil {
			// 区分不同请求方式存储api
			for _, api := range apis { // nolint:gocritic
				arr := make([]models.SysApi, 0)
				if list, ok := apiMap[api.Method]; ok {
					arr = append(list, api) // nolint:gocritic
				} else {
					arr = append(arr, api)
				}
				apiMap[api.Method] = arr
			}
		}
		// 写入缓存
		apiCache.Set("apiMap", apiMap, cache.DefaultExpiration)
	}

	// 匹配路由
	for _, api := range apiMap[method] { // nolint:gocritic
		// 通过casbin KeyMatch2来匹配url规则
		match := util.KeyMatch2(path, api.Path)
		if match {
			desc = api.Desc
			break
		}
	}

	// 写入缓存
	apiCache.Set(fmt.Sprintf("%s_%s", method, path), desc, cache.DefaultExpiration)

	return desc
}
