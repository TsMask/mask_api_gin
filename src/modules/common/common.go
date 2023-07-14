package common

import (
	"mask_api_gin/src/modules/common/controller"
	"mask_api_gin/src/pkg/logger"
	"mask_api_gin/src/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> common 模块路由")

	// 路由主页
	router.GET("/", controller.Index.Handler)

	// 验证码操作处理
	router.GET("/captchaImage", controller.Captcha.Image)

	// 账号身份操作处理
	router.POST("/login", controller.Account.Login)
	router.GET("/getInfo", middleware.PreAuthorize(nil), controller.Account.Info)
	router.GET("/getRouters", controller.Account.Router)
	router.POST("/logout", controller.Account.Logout)

	// 通用请求
	commonGroup := router.Group("/common")
	{
		// 路由主页
		commonGroup.GET("/hash", controller.Commont.Hash)
	}

	// 文件操作处理
	fileGroup := router.Group("/file")
	{
		// 下载文件
		fileGroup.GET("/download/:filePath", controller.File.Download)
		// 上传文件
		fileGroup.GET("/upload", controller.File.Upload)
		// 切片文件检查
		fileGroup.POST("/chunkCheck", controller.File.ChunkCheck)
		// 切片文件上传
		fileGroup.GET("/chunkUpload", controller.File.ChunkUpload)
		// 切片文件合并
		fileGroup.GET("/chunkMerge", controller.File.ChunkMerge)
	}
}
