package engine

import (
	"mask_api_gin/src/framework/catch"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/framework/middleware/security"

	"fmt"

	"github.com/gin-gonic/gin"
)

// NewApp 初始应用引擎
func NewApp() *gin.Engine {
	app := initAppEngine()

	initMiddleware(app)
	initStaticFS(app)

	// 路由未找到时
	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  fmt.Sprintf("Not Found %s %s", c.Request.Method, c.Request.RequestURI),
		})
	})

	return app
}

// initAppEngine 初始应用引擎
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

// initMiddleware 初始全局中间件
func initMiddleware(app *gin.Engine) {
	// 全局中间件
	app.Use(catch.ErrorCatch())
	if config.Env() == "local" {
		app.Use(middleware.Report())
	}
	app.Use(middleware.Cors(), security.Security())
}

// initStaticFS 初始静态文件路由
func initStaticFS(app *gin.Engine) {
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
}
