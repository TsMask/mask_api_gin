package controller

import (
	"fmt"
	"mask_api_gin/src/framework/model/result"
	"mask_api_gin/src/framework/service/ctx"

	"github.com/gin-gonic/gin"
)

// 路由主页
var Index = new(indexController)

type indexController struct{}

// 根路由
//
// GET /
func (s *indexController) Handler(c *gin.Context) {
	name := ctx.Config("pkg.name").(string)
	version := ctx.Config("pkg.version").(string)
	str := "欢迎使用%s后台管理框架，当前版本：%s，请通过前端管理地址访问。"
	c.JSON(200, result.OkMsg(fmt.Sprintf(str, name, version)))
}
