package common

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/modules/common/controller"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> common 模块路由")

	// 路由主页
	router.GET("/",
		middleware.RateLimit(middleware.LimitOption{
			Time:  300,
			Count: 10,
			Type:  middleware.LIMIT_IP,
		}),
		controller.NewIndex.Handler,
	)

	// 通用请求
	commonGroup := router.Group("/common")
	{
		commonGroup.POST("/hash", middleware.PreAuthorize(nil), controller.NewCommon.Hash)
	}

	// 文件操作
	fileGroup := router.Group("/file")
	{
		fileGroup.POST("/upload", middleware.PreAuthorize(nil), controller.NewFile.Upload)
		fileGroup.POST("/chunk-check", middleware.PreAuthorize(nil), controller.NewFile.ChunkCheck)
		fileGroup.POST("/chunk-upload", middleware.PreAuthorize(nil), controller.NewFile.ChunkUpload)
		fileGroup.POST("/chunk-merge", middleware.PreAuthorize(nil), controller.NewFile.ChunkMerge)
		fileGroup.GET("/download/:filePath", middleware.PreAuthorize(nil), controller.NewFile.Download)
	}
}
