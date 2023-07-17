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
	{
		indexGroup.GET("/",
			middleware.RateLimit(map[string]int64{
				"time":  300,
				"count": 10,
				"type":  middleware.LIMIT_IP,
			}),
			controller.Index.Handler,
		)

		// 验证码操作处理
		indexGroup.GET("/captchaImage",
			middleware.RateLimit(map[string]int64{
				"time":  300,
				"count": 60,
				"type":  middleware.LIMIT_IP,
			}),
			controller.Captcha.Image,
		)

		// 账号身份操作处理
		indexGroup.POST("/login",
			middleware.RateLimit(map[string]int64{
				"time":  300,
				"count": 10,
				"type":  middleware.LIMIT_IP,
			}),
			controller.Account.Login,
		)
		indexGroup.GET("/getInfo", middleware.PreAuthorize(nil), controller.Account.Info)
		indexGroup.GET("/getRouters", middleware.PreAuthorize(nil), controller.Account.Router)
		indexGroup.POST("/logout",
			middleware.RateLimit(map[string]int64{
				"time":  300,
				"count": 5,
				"type":  middleware.LIMIT_IP,
			}),
			controller.Account.Logout,
		)

		// 账号注册操作处理
		indexGroup.POST("/register",
			middleware.RateLimit(map[string]int64{
				"time":  300,
				"count": 10,
				"type":  middleware.LIMIT_IP,
			}),
			controller.Register.UserName,
		)
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
