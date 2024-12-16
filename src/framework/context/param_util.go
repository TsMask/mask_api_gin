package context

import (
	"mask_api_gin/src/framework/ip2region"
	"mask_api_gin/src/framework/utils/ua"

	"github.com/gin-gonic/gin"
)

// IPAddrLocation 解析ip地址
func IPAddrLocation(c *gin.Context) (string, string) {
	ip := ip2region.ClientIP(c.ClientIP())
	location := ip2region.RealAddressByIp(ip)
	return ip, location
}

// UaOsBrowser 解析请求用户代理信息
func UaOsBrowser(c *gin.Context) (string, string) {
	userAgent := c.GetHeader("user-agent")
	uaInfo := ua.Info(userAgent)

	browser := "-"
	if bName, bVersion := uaInfo.Browser(); bName != "" {
		browser = bName
		if bVersion != "" {
			browser = bName + " " + bVersion
		}
	}

	os := "-"
	if bos := uaInfo.OS(); bos != "" {
		os = bos
	}
	return os, browser
}
