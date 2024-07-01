package error_catch

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/vo/result"

	"github.com/gin-gonic/gin"
)

// ErrorCatch 全局异常捕获
func ErrorCatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// 在这里处理 Panic 异常，例如记录日志或返回错误信息给客户端
			if err := recover(); err != nil {
				logger.Errorf("Panic 异常: %s => %v", c.Request.URL, err)

				// 返回错误响应给客户端
				if config.Env() == "prod" {
					c.JSON(500, result.ErrMsg("服务器内部错误"))
				} else {
					// 通过实现 error 接口的 Error() 方法自定义错误类型进行捕获
					switch v := err.(type) {
					case error:
						c.JSON(500, result.ErrMsg(v.Error()))
					default:
						c.JSON(500, result.ErrMsg(fmt.Sprint(err)))
					}
				}

				c.Abort() // 停止执行后续的处理函数
			}
		}()

		c.Next()
	}
}
