package engine

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/web/middleware"
	"mask_api_gin/src/internal/web/middleware/security"
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

// initAppEngine 初始应用引擎，生产环境关闭 Gin 调试输出
func initAppEngine() *gin.Engine {
	gin.DisableConsoleColor()

	app := gin.New()
	if config.Env() == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		app.Use(gin.Logger())
	}
	return app
}

// initMiddleware 初始全局中间件
func initMiddleware(app *gin.Engine) {
	app.Use(middleware.Recovery())
	if config.Env() == "local" {
		app.Use(middleware.AccessLog())
	}
	app.Use(middleware.Cors(), security.Security())
}

// initStaticFS 初始静态文件路由
func initStaticFS(app *gin.Engine) {
	// 静态目录-静态资源
	prefix := config.GetStringOr("staticFile.default.prefix", "")
	dir := config.GetStringOr("staticFile.default.dir", "")
	if prefix != "" && dir != "" {
		app.StaticFS(prefix, gin.Dir(dir, true))
	}

	// 静态目录-上传资源
	prefix = config.GetStringOr("staticFile.upload.prefix", "")
	dir = config.GetStringOr("staticFile.upload.dir", "")
	if prefix != "" && dir != "" {
		app.StaticFS(prefix, gin.Dir(dir, true))
	}
}
