package src

import (
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/cron"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/redis"
)

// 配置中心初始加载
func ConfigurationInit() {
	// 初始配置参数
	config.InitConfig()
	// 连接数据库实例
	datasource.Connect()
	// 连接Redis实例
	redis.Connect()
	// 启动调度任务实例
	cron.StartCron()
}

// 配置中心相关配置关闭连接
func ConfigurationClose() {
	// 关闭数据库实例
	datasource.Close()
	// 关闭Redis实例
	redis.Close()
	// 停止调度任务实例
	cron.StopCron()
}
