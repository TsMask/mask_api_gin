package controller

import (
	"mask_api_gin/src/framework/vo/result"

	"github.com/gin-gonic/gin"
)

// 服务器监控
var Server = new(server)

type server struct{}

// 服务器信息
//
// GET /
func (s *server) Info(c *gin.Context) {
	c.JSON(200, result.OkMsg("sdnfo"))
}
