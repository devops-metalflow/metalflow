package middleware

import (
	"metalflow/pkg/global"
	"metalflow/pkg/response"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
)

// 幂等性中间件

var (
	expire = 24 * time.Hour
	// 记录是否加锁
	idempotenceLock sync.Mutex
	// 存取token
	idempotenceMap = cache.New(expire, 48*time.Hour) //nolint:gomnd
)

// Idempotence 全局异常处理中间件
func Idempotence(c *gin.Context) {
	// 优先从header提取
	token := c.Request.Header.Get(global.Conf.System.IdempotenceTokenName)
	if token == "" {
		token, _ = c.Cookie(global.Conf.System.IdempotenceTokenName)
	}
	token = strings.TrimSpace(token)
	if token == "" {
		response.FailWithMsg(response.IdempotenceTokenEmptyMsg)
	}
	// token校验
	if !CheckIdempotenceToken(token) {
		response.FailWithMsg(response.IdempotenceTokenInvalidMsg)
	}
	c.Next()
}

// GetIdempotenceToken 全局异常处理中间件
func GetIdempotenceToken(_ *gin.Context) {
	response.SuccessWithData(GenIdempotenceToken())
}

// GenIdempotenceToken 生成一个幂等性token
func GenIdempotenceToken() string {
	token := uuid.NewV4().String()
	// 写入map
	idempotenceLock.Lock()
	defer idempotenceLock.Unlock()
	idempotenceMap.Set(token, 1, cache.DefaultExpiration)
	return token
}

// CheckIdempotenceToken 校验幂等性token
func CheckIdempotenceToken(token string) bool {
	idempotenceLock.Lock()
	defer idempotenceLock.Unlock()
	// 读取map
	_, ok := idempotenceMap.Get(token)
	if !ok {
		return false
	}
	// 删除map中的值
	idempotenceMap.Delete(token)
	return true
}
