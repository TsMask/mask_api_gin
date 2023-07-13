package ctx

import (
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/utils/ip2region"
	"mask_api_gin/src/framework/utils/ua"
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

// 解析ip地址
func ClientIP(c *gin.Context) (string, string) {
	ip := c.ClientIP()
	ipAddr := ip
	location := common.IP_INNER_LOCATION
	if strings.Contains(ip, common.IP_INNER_ADDR) {
		ipAddr = strings.Replace(ip, common.IP_INNER_ADDR, "", -1)
	} else {
		location = ip2region.RealAddressByIp(ip)
	}
	return ipAddr, location
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
