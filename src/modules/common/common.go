package common

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/modules/common/controller"

	"github.com/gin-gonic/gin"
)

// 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> common 模块路由")

	// 路由主页
	indexGroup := router.Group("/")
	indexGroup.GET("",
		middleware.RateLimit(middleware.LimitOption{
			Time:  300,
			Count: 10,
			Type:  middleware.LIMIT_IP,
		}),
		controller.NewIndex.Handler,
	)

	// 验证码操作处理
	indexGroup.GET("/captchaImage",
		middleware.RateLimit(middleware.LimitOption{
			Time:  300,
			Count: 60,
			Type:  middleware.LIMIT_IP,
		}),
		controller.NewCaptcha.Image,
	)

	// 账号身份操作处理
	{
		indexGroup.POST("/login",
			middleware.RateLimit(middleware.LimitOption{
				Time:  300,
				Count: 10,
				Type:  middleware.LIMIT_IP,
			}),
			controller.NewAccount.Login,
		)
		indexGroup.GET("/getInfo", middleware.PreAuthorize(nil), controller.NewAccount.Info)
		indexGroup.GET("/getRouters", middleware.PreAuthorize(nil), controller.NewAccount.Router)
		indexGroup.POST("/logout",
			middleware.RateLimit(middleware.LimitOption{
				Time:  300,
				Count: 5,
				Type:  middleware.LIMIT_IP,
			}),
			controller.NewAccount.Logout,
		)
	}

	// 账号注册操作处理
	{
		indexGroup.POST("/register",
			middleware.RateLimit(middleware.LimitOption{
				Time:  300,
				Count: 10,
				Type:  middleware.LIMIT_IP,
			}),
			controller.NewRegister.UserName,
		)
	}

	// 通用请求
	commonGroup := router.Group("/common")
	{
		commonGroup.GET("/hash", middleware.PreAuthorize(nil), controller.NewCommont.Hash)
	}

	// 文件操作处理
	fileGroup := router.Group("/file")
	{
		fileGroup.GET("/download/:filePath", middleware.PreAuthorize(nil), controller.NewFile.Download)
		fileGroup.POST("/upload", middleware.PreAuthorize(nil), controller.NewFile.Upload)
		fileGroup.POST("/chunkCheck", middleware.PreAuthorize(nil), controller.NewFile.ChunkCheck)
		fileGroup.POST("/chunkUpload", middleware.PreAuthorize(nil), controller.NewFile.ChunkUpload)
		fileGroup.POST("/chunkMerge", middleware.PreAuthorize(nil), controller.NewFile.ChunkMerge)
	}
}
