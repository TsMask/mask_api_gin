package controller

import (
	"github.com/gin-gonic/gin"
)

// 实例化控制层 CommontController 结构体
var NewCommont = &CommontController{}

// 通用请求
//
// PATH /
type CommontController struct{}

// 哈希加密
//
// GET /hash
func (s *CommontController) Hash(c *gin.Context) {
	c.String(200, "commont Hash")
}
