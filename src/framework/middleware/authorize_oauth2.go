package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/token"
)

// AuthorizeOauth2 客户端授权认证校验
//
// scope 客户端授权范围，例如：[]string{"read","write"}
func AuthorizeOauth2(scope []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头标识信息
		tokenStr := reqctx.Authorization(c)
		if tokenStr == "" {
			c.JSON(401, resp.CodeMsg(401, "无效客户端授权"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 验证令牌
		claims, err := token.Oauth2TokenVerify(tokenStr)
		if err != nil {
			c.JSON(401, resp.CodeMsg(40003, err.Error()))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 获取缓存的用户信息
		info := token.Oauth2InfoGet(claims)
		if info.ClientId == "" {
			c.JSON(401, resp.CodeMsg(401, "无效客户端授权"))
			c.Abort() // 停止执行后续的处理函数
			return
		}
		c.Set(constants.CTX_LOGIN_OAUTH2, info)

		// 客户端权限校验
		if scope != nil {
			var hasScope bool = false
			for _, item := range info.Scope {
				for _, v := range scope {
					if item == v {
						hasScope = true
						break
					}
				}
			}
			if !hasScope {
				msg := fmt.Sprintf("无权访问 %s %s", c.Request.Method, c.Request.RequestURI)
				c.JSON(403, resp.CodeMsg(403, msg))
				c.Abort() // 停止执行后续的处理函数
				return
			}
		}

		// 调用下一个处理程序
		c.Next()
	}
}
