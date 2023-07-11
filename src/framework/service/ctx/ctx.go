package ctx

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 查询参数转换MapString
func QueryMapString(c *gin.Context) map[string]string {
	queryValues := c.Request.URL.Query()
	queryParams := make(map[string]string)
	for key, values := range queryValues {
		queryParams[key] = values[0]
	}
	return queryParams
}

// 获取运行服务环境
// local prod
func Env() string {
	return viper.GetString("env")
}

// 获取配置信息
//
// ctx.Config("pkg.name")
func Config(key string) interface{} {
	return viper.Get(key)
}
