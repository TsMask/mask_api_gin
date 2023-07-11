package src

import (
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/modules/common"
	"mask_api_gin/src/modules/monitor"
	"mask_api_gin/src/modules/system"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 初始应用引擎
func initAppEngine() *gin.Engine {
	var app *gin.Engine

	// 禁止控制台日志输出的颜色
	gin.DisableConsoleColor()
	// 根据运行环境注册引擎
	if viper.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
		app = gin.New()
		app.Use(gin.Recovery())
	} else {
		app = gin.Default()
	}

	// 全局中间件
	app.Use(middleware.LoggerMiddleware())

	// 静态目录
	fsDefault := viper.GetStringMapString("staticFile.default")
	app.StaticFS(fsDefault["prefix"], gin.Dir(fsDefault["dir"], true))
	fsUpload := viper.GetStringMapString("staticFile.upload")
	app.StaticFS(fsUpload["prefix"], gin.Dir(fsUpload["dir"], true))

	// 测试启动
	app.GET("/ping", func(c *gin.Context) {
		forwardedFor := c.Request.Header.Get("X-Forwarded-For")
		ip := c.ClientIP()

		c.JSON(200, gin.H{
			"config":       viper.AllSettings(),
			"forwarded_ip": forwardedFor,
			"client_ip":    ip,
		})
	})

	return app
}

// 初始模块路由
func initModulesRoute(app *gin.Engine) {
	common.Setup(app)
	monitor.Setup(app)
	system.Setup(app)
}

// 运行服务程序
func RunServer() error {
	app := initAppEngine()
	initModulesRoute(app)

	// 读取服务配置
	app.ForwardedByClientIP = viper.GetBool("server.proxy")
	addr := ":" + viper.GetString("server.port")

	// 启动服务
	return app.Run(addr)
}
