package reqctx

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/infra/ip2region"
	"mask_api_gin/src/internal/core/utils/crypto"
	"mask_api_gin/src/internal/core/utils/ua"
)

// IPAddrLocation 解析ip地址
func IPAddrLocation(c *gin.Context) (string, string) {
	ip := ip2region.ClientIP(c.ClientIP())
	location := ip2region.RealAddressByIp(ip)
	return ip, location
}

// UaOsBrowser 解析请求用户代理信息
func UaOsBrowser(c *gin.Context) (string, string) {
	userAgent := c.Request.UserAgent()
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

// DeviceFingerprint 设备指纹信息
func DeviceFingerprint(c *gin.Context, v any) string {
	str := fmt.Sprintf("%v:%s", v, c.Request.UserAgent())
	return crypto.SHA256ToBase64(str)
}
