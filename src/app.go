package src

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	errorCatch "mask_api_gin/src/framework/error_catch"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/framework/middleware/security"
	"mask_api_gin/src/modules/common"
	"mask_api_gin/src/modules/demo"
	"mask_api_gin/src/modules/monitor"
	"mask_api_gin/src/modules/system"

	"github.com/gin-gonic/gin"
)

// RunServer 运行服务程序
func RunServer() error {
	app := initAppEngine()

	// 初始全局默认
	initDefeat(app)

	// 初始模块路由
	initModulesRoute(app)

	// 读取服务配置
	app.ForwardedByClientIP = config.Get("server.proxy").(bool)
	addr := fmt.Sprintf(":%d", config.Get("server.port").(int))

	// 启动服务
	fmt.Printf("\nopen http://localhost%s \n\n", addr)
	return app.Run(addr)
}

// 初始应用引擎
func initAppEngine() *gin.Engine {
	var app *gin.Engine

	// 禁止控制台日志输出的颜色
	gin.DisableConsoleColor()

	// 根据运行环境注册引擎
	if config.Env() == "prod" {
		gin.SetMode(gin.ReleaseMode)
		app = gin.New()
		app.Use(gin.Recovery())
	} else {
		app = gin.Default()
	}

	return app
}

// 初始全局默认
func initDefeat(app *gin.Engine) {
	// 全局中间件
	app.Use(errorCatch.ErrorCatch())
	if config.Env() == "local" {
		app.Use(middleware.Report())
	}
	app.Use(middleware.Cors(), security.Security())

	// 静态目录-静态资源
	if v := config.Get("staticFile.default"); v != nil {
		fsMap := v.(map[string]any)
		prefix, dir := fsMap["prefix"], fsMap["dir"]
		if prefix != nil && dir != nil {
			app.StaticFS(prefix.(string), gin.Dir(dir.(string), true))
		}
	}

	// 静态目录-上传资源
	if v := config.Get("staticFile.upload"); v != nil {
		fsMap := v.(map[string]any)
		prefix, dir := fsMap["prefix"], fsMap["dir"]
		if prefix != nil && dir != nil {
			app.StaticFS(prefix.(string), gin.Dir(dir.(string), true))
		}
	}

	// 路由未找到时
	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  fmt.Sprintf("Not Found %s %s", c.Request.Method, c.Request.RequestURI),
		})
	})
}

// 初始模块路由
func initModulesRoute(app *gin.Engine) {
	demo.Setup(app) // 演示模块（可移除）

	common.Setup(app)  // 通用模块
	system.Setup(app)  // 系统模块
	monitor.Setup(app) // 监控模块（含定时任务放最后）
}
