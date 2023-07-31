package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/vo/result"

	"github.com/gin-gonic/gin"
)

// 实例化控制层 IndexController 结构体
var NewIndex = &IndexController{}

// 路由主页
//
// PATH /
type IndexController struct{}

// 根路由
//
// GET /
func (s *IndexController) Handler(c *gin.Context) {
	name := config.Get("framework.name").(string)
	version := config.Get("framework.version").(string)
	str := "欢迎使用%s后台管理框架，当前版本：%s，请通过前端管理地址访问。"
	c.JSON(200, result.OkMsg(fmt.Sprintf(str, name, version)))
}
