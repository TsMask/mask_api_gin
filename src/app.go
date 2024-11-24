package src

import (
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/engine"
	"mask_api_gin/src/modules/common"
	"mask_api_gin/src/modules/demo"
	"mask_api_gin/src/modules/monitor"
	"mask_api_gin/src/modules/system"

	"fmt"

	"github.com/gin-gonic/gin"
)

// RunServer 运行服务程序
func RunServer() error {
	app := engine.NewApp()

	// 模块路由注册
	ModulesSetup(app)

	// 读取服务配置
	app.ForwardedByClientIP = config.Get("server.proxy").(bool)
	addr := fmt.Sprintf("%s:%d", config.Get("server.host"), config.Get("server.port"))

	// 启动服务
	fmt.Printf("\nopen http://%s \n\n", addr)
	return app.Run(addr)
}

// 模块路由注册
func ModulesSetup(app *gin.Engine) {
	demo.Setup(app) // 演示模块（可移除）

	common.Setup(app)  // 通用模块
	system.Setup(app)  // 系统模块
	monitor.Setup(app) // 监控模块（含定时任务放最后）
}
