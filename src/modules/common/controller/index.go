package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/vo/result"

	"github.com/gin-gonic/gin"
)

// 路由主页
var Index = new(indexController)

type indexController struct{}

// 根路由
//
// GET /
func (s *indexController) Handler(c *gin.Context) {
	name := config.Get("framework.name").(string)
	version := config.Get("framework.version").(string)
	str := "欢迎使用%s后台管理框架，当前版本：%s，请通过前端管理地址访问。"
	c.JSON(200, result.OkMsg(fmt.Sprintf(str, name, version)))
}
