package common

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/core/logger"
	"mask_api_gin/src/internal/server/middleware"
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
		commonGroup.POST("/hash", middleware.AuthorizeUser(nil), controller.NewCommon.Hash)
	}

	// 文件操作
	fileGroup := router.Group("/file")
	file := controller.NewFile
	{
		fileGroup.POST("/upload", middleware.AuthorizeUser(nil), file.Upload)
		fileGroup.POST("/chunk-check", middleware.AuthorizeUser(nil), file.ChunkCheck)
		fileGroup.POST("/chunk-upload", middleware.AuthorizeUser(nil), file.ChunkUpload)
		fileGroup.POST("/chunk-merge", middleware.AuthorizeUser(nil), file.ChunkMerge)
		fileGroup.GET("/download/:filePath", middleware.AuthorizeUser(nil), file.Download)
	}
}
