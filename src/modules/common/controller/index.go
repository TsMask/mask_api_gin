package controller

import (
	"fmt"
	"mask_api_gin/src/framework/model/result"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 路由主页
var Index = new(index)

type index struct{}

// 根路由
//
// GET /
func (s *index) Handler(c *gin.Context) {
	name := viper.GetString("pkg.name")
	version := viper.GetString("pkg.version")
	str := "欢迎使用%s后台管理框架，当前版本：%s，请通过前端管理地址访问。"
	c.JSON(200, result.OkMsg(fmt.Sprintf(str, name, version)))
}
