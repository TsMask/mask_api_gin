package ctx

import (
	constToken "mask_api_gin/src/framework/constants/token"

	"strings"

	"github.com/gin-gonic/gin"
)

// QueryMap Query参数转换Map
func QueryMap(c *gin.Context) map[string]any {
	queryValues := c.Request.URL.Query()
	queryParams := make(map[string]any, len(queryValues))
	for key, values := range queryValues {
		queryParams[key] = values[0]
	}
	return queryParams
}

// BodyJSONMap JSON参数转换Map
func BodyJSONMap(c *gin.Context) map[string]any {
	params := make(map[string]any, 0)
	c.ShouldBindBodyWithJSON(&params)
	return params
}

// RequestParamsMap 请求参数转换Map
func RequestParamsMap(c *gin.Context) map[string]any {
	params := make(map[string]any, 0)
	// json
	if strings.HasPrefix(c.ContentType(), "application/json") {
		c.ShouldBindBodyWithJSON(&params)
	}

	// 表单
	formParams := c.Request.PostForm
	for key, value := range formParams {
		if _, ok := params[key]; !ok {
			params[key] = value[0]
		}
	}

	// 查询
	queryParams := c.Request.URL.Query()
	for key, value := range queryParams {
		if _, ok := params[key]; !ok {
			params[key] = value[0]
		}
	}
	return params
}

// Authorization 解析请求头
func Authorization(c *gin.Context) string {
	authHeader := c.GetHeader(constToken.HEADER_KEY)
	if authHeader == "" {
		return ""
	}
	// 拆分 Authorization 请求头，提取 JWT 令牌部分
	tokenStr := strings.TrimPrefix(authHeader, constToken.HEADER_PREFIX)
	if len(tokenStr) > 64 {
		return strings.TrimSpace(tokenStr) // 去除可能存在的空格
	}
	return ""
}
