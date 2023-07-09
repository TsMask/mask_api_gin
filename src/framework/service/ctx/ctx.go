package ctx

import "github.com/gin-gonic/gin"

// 查询参数转换MapString
func QueryMapString(c *gin.Context) map[string]string {
	queryValues := c.Request.URL.Query()
	queryParams := make(map[string]string)
	for key, values := range queryValues {
		queryParams[key] = values[0]
	}
	return queryParams
}
