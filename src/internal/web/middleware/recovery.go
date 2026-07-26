package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
)

// Recovery 捕获处理链中的 panic，记录日志并返回统一错误响应，
// 避免单个请求的异常导致整个进程崩溃。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				userId := reqctx.LoginUserToUserID(c)
				logger.Errorf("[Recovery][%d] panic recovered: %s => %v", userId, c.Request.URL, err)

				// 生产环境隐藏内部错误细节
				if config.Env() == "prod" {
					c.JSON(500, resp.CodeMsg(resp.CODE_INTERNAL, resp.MSG_INTERNAL))
				} else {
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
