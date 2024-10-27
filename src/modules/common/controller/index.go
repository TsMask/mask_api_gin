package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/vo/result"

	"github.com/gin-gonic/gin"
)

// NewIndex 实例化控制层
var NewIndex = &IndexController{}

// IndexController 路由主页 控制层处理
//
// PATH /
type IndexController struct{}

// Handler 根路由
//
// GET /
func (s *IndexController) Handler(c *gin.Context) {
	str := "欢迎使用%s后台管理框架，当前版本：%s，请通过前端管理地址访问。"
	c.JSON(200, result.OkMsg(fmt.Sprintf(str, config.Name, config.Version)))
}
