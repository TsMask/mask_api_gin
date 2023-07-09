package controller

import (
	"mask_api_gin/src/framework/model/result"

	"github.com/gin-gonic/gin"
)

// 参数配置信息
var SysConfig = new(sys_config)

type sys_config struct{}

// 导出参数配置信息
//
// POST /export
func (s *sys_config) Export(c *gin.Context) {
	c.JSON(200, result.OkMsg("sdnfo"))
}
