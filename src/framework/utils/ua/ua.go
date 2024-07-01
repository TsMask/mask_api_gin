package ua

import "github.com/mssola/user_agent"

// Info 获取user-agent信息
func Info(userAgent string) *user_agent.UserAgent {
	return user_agent.New(userAgent)
}
