package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
)

// ErrorCatch 全局异常捕获
func ErrorCatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// 在这里处理 Panic 异常，例如记录日志或返回错误信息给客户端
			if err := recover(); err != nil {
				userId := reqctx.LoginUserToUserID(c)
				logger.Errorf("[Recovery][%d] panic recovered: %s => %v", userId, c.Request.URL, err)

				// 返回错误响应给客户端
				if config.Env() == "prod" {
					c.JSON(500, resp.CodeMsg(resp.CODE_INTERNAL, resp.MSG_INTERNAL))
				} else {
					// 通过实现 error 接口的 Error() 方法自定义错误类型进行捕获
					switch v := err.(type) {
					case error:
						c.JSON(500, resp.CodeMsg(resp.CODE_INTERNAL, v.Error()))
					default:
						c.JSON(500, resp.CodeMsg(resp.CODE_INTERNAL, fmt.Sprint(err)))
					}
				}

				c.Abort() // 停止执行后续的处理函数
			}
		}()

		c.Next()
	}
}
