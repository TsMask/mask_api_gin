package common

import (
	"mask_api_gin/src/modules/common/controller"

	"github.com/gin-gonic/gin"
)

// 模块路由注册
func Setup(router *gin.Engine) {
	// 根路由组
	indexGroup := router.Group("/")
	{
		// 路由主页
		indexGroup.GET("/", controller.Index.Handler)
		// 获取验证码
		indexGroup.GET("/captchaImage", controller.Captcha.CaptchaImage)
	}

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
