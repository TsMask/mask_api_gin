package common

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/web/middleware"
	"mask_api_gin/src/modules/common/controller"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> common 模块路由")

	// 路由主页
	router.GET("/",
		middleware.RateLimit(300, 10, middleware.RateLimitIP),
		controller.NewIndex.Handler,
	)

	// 健康检查：存活无状态，就绪检查依赖可用性。
	router.GET("/healthz", controller.NewHealth.Live)
	router.GET("/readyz", controller.NewHealth.Ready)

	// 通用请求
	commonGroup := router.Group("/api/v1/common")
	{
		commonGroup.POST("/hash", middleware.AuthorizeLogin(), controller.NewCommon.Hash)
	}

	// 文件操作
	fileGroup := router.Group("/api/v1/files")
	file := controller.NewFile
	{
		fileGroup.POST("/upload", middleware.AuthorizeLogin(), file.Upload)
		fileGroup.POST("/chunk-check", middleware.AuthorizeLogin(), file.ChunkCheck)
		fileGroup.POST("/chunk-upload", middleware.AuthorizeLogin(), file.ChunkUpload)
		fileGroup.POST("/chunk-merge", middleware.AuthorizeLogin(), file.ChunkMerge)
		fileGroup.GET("/download", middleware.AuthorizeLogin(), file.Download)
	}
}
