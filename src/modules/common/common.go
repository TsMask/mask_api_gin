package common

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/modules/common/controller"

	"github.com/gin-gonic/gin"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> common 模块路由")

	// 路由主页
	router.GET("/",
		middleware.RateLimit(middleware.LimitOption{
			Time:  300,
			Count: 10,
			Type:  middleware.LimitIP,
		}),
		controller.NewIndex.Handler,
	)

	// 验证码操作
	router.GET("/captchaImage",
		middleware.RateLimit(middleware.LimitOption{
			Time:  300,
			Count: 60,
			Type:  middleware.LimitIP,
		}),
		controller.NewCaptcha.Image,
	)

	// 账号身份操作
	{
		router.POST("/login",
			middleware.RateLimit(middleware.LimitOption{
				Time:  180,
				Count: 15,
				Type:  middleware.LimitIP,
			}),
			controller.NewAccount.Login,
		)
		router.GET("/getInfo", middleware.PreAuthorize(nil), controller.NewAccount.Info)
		router.GET("/getRouters", middleware.PreAuthorize(nil), controller.NewAccount.Router)
		router.POST("/logout",
			middleware.RateLimit(middleware.LimitOption{
				Time:  120,
				Count: 15,
				Type:  middleware.LimitIP,
			}),
			controller.NewAccount.Logout,
		)
	}

	// 账号注册操作
	{
		router.POST("/register",
			middleware.RateLimit(middleware.LimitOption{
				Time:  300,
				Count: 10,
				Type:  middleware.LimitIP,
			}),
			controller.NewRegister.Register,
		)
	}

	// 通用请求
	commonGroup := router.Group("/common")
	{
		commonGroup.POST("/hash", middleware.PreAuthorize(nil), controller.NewCommon.Hash)
	}

	// 文件操作
	fileGroup := router.Group("/file")
	{
		fileGroup.GET("/download/:filePath", middleware.PreAuthorize(nil), controller.NewFile.Download)
		fileGroup.POST("/upload", middleware.PreAuthorize(nil), controller.NewFile.Upload)
		fileGroup.POST("/chunkCheck", middleware.PreAuthorize(nil), controller.NewFile.ChunkCheck)
		fileGroup.POST("/chunkUpload", middleware.PreAuthorize(nil), controller.NewFile.ChunkUpload)
		fileGroup.POST("/chunkMerge", middleware.PreAuthorize(nil), controller.NewFile.ChunkMerge)
	}
}
