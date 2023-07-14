package src

import (
	"mask_api_gin/src/pkg/cache/redis"
	"mask_api_gin/src/pkg/config"
	"mask_api_gin/src/pkg/datasource"
)

// 配置中心初始加载
func ConfigurationInit() {
	// 初始配置参数
	config.InitConfig("./src/config")
	// 连接数据库实例
	datasource.Connect()
	// 连接Redis实例
	redis.Connect()
}

// 配置中心相关配置关闭连接
func ConfigurationClose() {
	// 关闭数据库实例
	datasource.Close()
	// 关闭Redis实例
	redis.Close()
}
