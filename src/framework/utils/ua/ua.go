package ua

import "github.com/mssola/useragent"

// Info 获取user-agent信息
func Info(userAgent string) *useragent.UserAgent {
	return useragent.New(userAgent)
}
