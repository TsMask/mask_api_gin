package controller

import (
	"github.com/gin-gonic/gin"
)

// NewCommon 实例化控制层
var NewCommon = &CommonController{}

// CommonController 通用请求 控制层处理
//
// PATH /
type CommonController struct{}

// Hash 哈希加密
//
// GET /hash
func (s *CommonController) Hash(c *gin.Context) {
	c.String(200, "Common Hash")
}
