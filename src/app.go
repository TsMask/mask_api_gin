package src

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/server/engine"
	"mask_api_gin/src/modules/auth"
	"mask_api_gin/src/modules/common"
	"mask_api_gin/src/modules/cronjob"
	"mask_api_gin/src/modules/demo"
	"mask_api_gin/src/modules/monitor"
	"mask_api_gin/src/modules/system"
)

// RunServer 运行服务程序
func RunServer() error {
	app := engine.NewApp()

	// 模块路由注册
	modulesSetup(app)

	// 读取服务配置
	app.ForwardedByClientIP = config.Get("server.proxy").(bool)
	addr := fmt.Sprintf("%s:%d", config.Get("server.host"), config.Get("server.port"))

	// 启动服务
	serverCert := fmt.Sprint(config.Get("server.cert"))
	serverkey := fmt.Sprint(config.Get("server.key"))
	if serverCert != "" && serverkey != "" && serverCert != "<nil>" && serverkey != "<nil>" {
		fmt.Printf("\nListening and serving HTTPS on %s \n\n", addr)
		return app.RunTLS(addr, serverCert, serverkey)
	}
	fmt.Printf("\nListening and serving HTTP on %s \n\n", addr)
	return app.Run(addr)
}

// 模块路由注册
func modulesSetup(app *gin.Engine) {
	demo.Setup(app) // 演示模块（可移除）

	auth.Setup(app)    // 认证模块
	common.Setup(app)  // 通用模块
	system.Setup(app)  // 系统模块
	cronjob.Setup()    // 任务模块
	monitor.Setup(app) // 监控模块（对任务名称绑定放最后）
}
