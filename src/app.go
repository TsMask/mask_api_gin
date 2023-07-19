package src

import (
	"fmt"
	"mask_api_gin/src/framework/config"
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
	if config.Env() == "prod" {
		gin.SetMode(gin.ReleaseMode)
		app = gin.New()
		app.Use(gin.Recovery())
	} else {
		app = gin.Default()
	}

	// 全局中间件
	app.Use(middleware.Report())

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

	// 路由未找到时
	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  c.Request.RequestURI + " Not Found",
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
	app.ForwardedByClientIP = config.Get("server.proxy").(bool)
	addr := fmt.Sprintf(":%d", config.Get("server.port").(int))

	// 启动服务
	return app.Run(addr)
}
