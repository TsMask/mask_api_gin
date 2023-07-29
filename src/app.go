package src

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/modules/common"
	"mask_api_gin/src/modules/demo"
	"mask_api_gin/src/modules/monitor"
	"mask_api_gin/src/modules/system"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 运行服务程序
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
	app.Use(middleware.Report(), middleware.Cors())

	// 静态目录
	fsDefault := viper.GetStringMapString("staticFile.default")
	app.StaticFS(fsDefault["prefix"], gin.Dir(fsDefault["dir"], true))
	fsUpload := viper.GetStringMapString("staticFile.upload")
	app.StaticFS(fsUpload["prefix"], gin.Dir(fsUpload["dir"], true))

	// 路由未找到时
	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  c.Request.RequestURI + " Not Found",
		})
	})
}

// 初始模块路由
func initModulesRoute(app *gin.Engine) {
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

	demo.Setup(app)
	common.Setup(app)
	monitor.Setup(app)
	system.Setup(app) // 一定放最后，定时任务加载
}
