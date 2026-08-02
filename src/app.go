package src

import (
	"fmt"
	"net"
	"strconv"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/web/engine"
	"mask_api_gin/src/modules/auth"
	"mask_api_gin/src/modules/common"
	"mask_api_gin/src/modules/demo"
	"mask_api_gin/src/modules/job"
	"mask_api_gin/src/modules/monitor"
	"mask_api_gin/src/modules/system"
)

// RunServer 运行服务程序（阻塞，直到监听出错才返回）。
func RunServer() error {
	app := engine.NewApp()

	// 模块路由注册
	modulesSetup(app)

	// 读取服务配置
	app.ForwardedByClientIP = config.GetBoolOr("server.proxy", false)
	host := config.GetStringOr("server.host", "")
	port := config.GetIntOr("server.port", 0)
	// net.JoinHostPort 自动给 IPv6 地址加方括号，避免 host:port 拼接歧义
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	// 启动服务
	serverCert := config.GetStringOr("server.cert", "")
	serverKey := config.GetStringOr("server.key", "")
	if serverCert != "" && serverKey != "" {
		fmt.Printf("\nListening and serving HTTPS on %s \n\n", addr)
		return app.RunTLS(addr, serverCert, serverKey)
	}
	fmt.Printf("\nListening and serving HTTP on %s \n\n", addr)
	return app.Run(addr)
}

// 模块路由注册
func modulesSetup(app *gin.Engine) {
	// 注入操作日志与任务执行日志持久化钩子（须在路由注册、任务初始化之前）
	system.RegisterOnOperateLog()
	job.RegisterOnJobLog()

	demo.Setup(app)    // 演示模块（可移除）
	auth.Setup(app)    // 认证模块
	common.Setup(app)  // 通用模块
	system.Setup(app)  // 系统模块
	job.Setup(app)     // 任务模块
	monitor.Setup(app) // 监控模块
}
