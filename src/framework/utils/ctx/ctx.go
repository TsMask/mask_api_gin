package ctx

import (
	"errors"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/constants/token"
	"mask_api_gin/src/framework/utils/ip2region"
	"mask_api_gin/src/framework/utils/ua"
	"mask_api_gin/src/framework/vo"
	"strings"

	"github.com/gin-gonic/gin"
)

// QueryMapString 查询参数转换MapString
func QueryMapString(c *gin.Context) map[string]string {
	queryValues := c.Request.URL.Query()
	queryParams := make(map[string]string)
	for key, values := range queryValues {
		queryParams[key] = values[0]
	}
	return queryParams
}

// IPAddrLocation 解析ip地址
func IPAddrLocation(c *gin.Context) (string, string) {
	ip := c.ClientIP()
	location := ip2region.RealAddressByIp(ip)
	return ip, location
}

// Authorization 解析请求头
func Authorization(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	// 拆分 Authorization 请求头，提取 JWT 令牌部分
	arr := strings.Split(authHeader, token.HEADER_PREFIX)
	if len(arr) == 2 && arr[1] == "" {
		return ""
	}
	return arr[1]
}

// UaOsBrowser 解析请求用户代理信息
func UaOsBrowser(c *gin.Context) (string, string) {
	userAgent := c.GetHeader("user-agent")
	uaInfo := ua.Info(userAgent)

	browser := "未知 未知"
	bName, bVersion := uaInfo.Browser()
	if bName != "" && bVersion != "" {
		browser = bName + " " + bVersion
	}

	os := "未知 未知"
	bos := uaInfo.OS()
	if bos != "" {
		os = bos
	}
	return os, browser
}

// LoginUser 登录用户信息
func LoginUser(c *gin.Context) (vo.LoginUser, error) {
	value, exists := c.Get(common.CTX_LOGIN_USER)
	if exists {
		return value.(vo.LoginUser), nil
	}
	return vo.LoginUser{}, errors.New("无效登录用户信息")
}
