package reqctx

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/core/constants"
)

// QueryMap Query参数转换Map
func QueryMap(c *gin.Context) map[string]string {
	queryValues := c.Request.URL.Query()
	queryParams := make(map[string]string, len(queryValues))
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

// Authorization 解析请求头携带的令牌
func Authorization(c *gin.Context) string {
	authHeader := c.GetHeader(constants.HEADER_KEY)
	if authHeader == "" {
		return ""
	}
	// 拆分 Authorization 请求头，提取 JWT 令牌部分
	tokenStr := strings.Replace(authHeader, constants.HEADER_PREFIX, "", 1)
	if len(tokenStr) > 64 {
		return strings.TrimSpace(tokenStr) // 去除可能存在的空格
	}
	return ""
}

// HeaderContentOctetStream 设置响应头为application/octet-stream返回二进制数据
func HeaderContentOctetStream(c *gin.Context, fileName string, fileData []byte) {
	baseName := filepath.Base(fileName)
	safeName := strings.ReplaceAll(baseName, "..", "_")
	safeName = strings.ReplaceAll(safeName, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	safeName = strings.ReplaceAll(safeName, ";", "_")
	safeName = strings.ReplaceAll(safeName, "\"", "_")

	escapedName := url.QueryEscape(safeName)
	disposition := fmt.Sprintf("attachment; filename=%s", escapedName)

	// 设置响应头
	c.Header("Content-Disposition", disposition)
	c.Header("Content-Type", "application/octet-stream")
	// 返回响应体
	c.Data(200, "application/octet-stream", fileData)
}
