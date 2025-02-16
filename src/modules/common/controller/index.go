package controller

import (
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
	c.JSON(200, resp.OkData(map[string]string{
		"name":    config.Name,
		"version": config.Version,
	}))
}
