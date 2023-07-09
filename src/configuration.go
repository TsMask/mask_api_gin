package src

import (
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/datasource"
)

// 运行配置中心
func Configuration() {
	// 初始配置参数
	config.InitConfig("./src/config")
	// 连接数据库实例
	datasource.Connect()
}
