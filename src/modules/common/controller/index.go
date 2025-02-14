package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/resp"
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
func (s IndexController) Handler(c *gin.Context) {
	msg := fmt.Sprintf("Name:%s, Version:%s", config.Name, config.Version)
	c.JSON(200, resp.OkMsg(msg))
}
