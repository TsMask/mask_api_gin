package src

import (
	"embed"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/cron"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/redis"
)

// ConfigurationInit 配置中心初始加载
func ConfigurationInit(assetsDir, configDir embed.FS) {
	// 初始配置参数
	config.InitConfig(configDir)
	config.SetAssetsDirFS(assetsDir)
	// 初始程序日志
	logger.InitLogger()
	// 连接数据库实例
	db.Connect()
	// 连接Redis实例
	redis.Connect()
	// 启动调度任务实例
	cron.StartCron()
}

// ConfigurationClose 配置中心相关配置关闭连接
func ConfigurationClose() {
	// 停止调度任务实例
	cron.StopCron()
	// 关闭Redis实例
	redis.Close()
	// 关闭数据库实例
	db.Close()
	// 关闭程序日志
	logger.Close()
}
