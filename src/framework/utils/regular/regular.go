package regular

import (
	"mask_api_gin/src/framework/constants/common"
	"regexp"
	"strings"
)

// Replace 正则替换
func Replace(originStr, regStr, repStr string) string {
	reg := regexp.MustCompile(regStr)
	return reg.ReplaceAllString(originStr, repStr)
}

// 判断是否为http(s)://开头
//
// link 网络链接
func ValidHttp(link string) bool {
	if link == "" {
		return false
	}
	return strings.HasPrefix(link, common.HTTP) || strings.HasPrefix(link, common.HTTPS)
}
