package middleware

import (
	"fmt"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"

	jwt "github.com/appleboy/gin-jwt/v2"
)

func InitAuth() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           global.Conf.Jwt.Realm,                                 // jwt标识
		Key:             []byte(global.Conf.Jwt.Key),                           // 服务端密钥
		Timeout:         time.Hour * time.Duration(global.Conf.Jwt.Timeout),    // token过期时间
		MaxRefresh:      time.Hour * time.Duration(global.Conf.Jwt.MaxRefresh), // token最大刷新时间(RefreshToken过期时间=Timeout+MaxRefresh)
		PayloadFunc:     payloadFunc,                                           // 有效载荷处理
		IdentityHandler: identityHandler,                                       // 解析Claims
		Authenticator:   login,                                                 // 校验token的正确性, 处理登录逻辑
		Authorizator:    authorizator,                                          // 用户登录校验成功处理
		Unauthorized:    unauthorized,                                          // 用户登录校验失败处理
		LoginResponse:   loginResponse,                                         // 登录成功后的响应
		LogoutResponse:  logoutResponse,                                        // 登出后的响应
		RefreshResponse: refreshResponse,                                       // 刷新token后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",    // 自动在这几个地方寻找请求中的token
		TokenHeadName:   "Bearer",                                              // header名称
		TimeFunc:        time.Now,
	})
}

func payloadFunc(data any) jwt.MapClaims {
	if v, ok := data.(map[string]any); ok {
		var user models.SysUser
		// 将用户json转为结构体
		utils.JsonI2Struct(v["user"], &user)
		return jwt.MapClaims{
			jwt.IdentityKey: user.Id,
			"user":          v["user"],
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) any {
	claims := jwt.ExtractClaims(c)
	// 此处返回值类型map[string]any与payloadFunc和authorizator的data类型必须一致, 否则会导致授权失败还不容易找到原因
	return map[string]any{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

func login(c *gin.Context) (any, error) {
	var req request.UserAuthRequestStruct
	// 请求json绑定
	_ = c.ShouldBindJSON(&req)

	// 创建服务
	s := service.New(c)
	// 密码校验
	user, err := s.LoginCheck(&req)
	if err != nil {
		return nil, err
	}
	// 将用户以json格式写入, payloadFunc/authorizator会使用到
	return map[string]any{
		"user": utils.Struct2Json(user),
	}, nil
}

func authorizator(data any, c *gin.Context) bool {
	if v, ok := data.(map[string]any); ok {
		userStr := v["user"].(string)
		var user models.SysUser
		utils.Json2Struct(userStr, &user)
		c.Set("user", user)
		return true
	}
	return false
}

func unauthorized(_ *gin.Context, code int, message string) {
	global.Log.Debug(fmt.Sprintf("JWT认证失败, 错误码%d, 错误信息：%s", code, message))
	if message == response.LoginCheckErrorMsg {
		response.FailWithMsg(response.LoginCheckErrorMsg)
		return
	}
	response.FailWithCode(response.Unauthorized)
}

func loginResponse(_ *gin.Context, _ int, token string, expires time.Time) {
	response.SuccessWithData(map[string]any{
		"token": token,
		"expires": models.LocalTime{
			Time: expires,
		},
	})
}

func logoutResponse(_ *gin.Context, _ int) {
	response.Success()
}

func refreshResponse(_ *gin.Context, _ int, token string, expires time.Time) {
	response.SuccessWithData(map[string]any{
		"token": token,
		"expires": models.LocalTime{
			Time: expires,
		},
	})
}
