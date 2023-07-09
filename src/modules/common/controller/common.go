package controller

import (
	"github.com/gin-gonic/gin"
)

// 通用请求
var Commont = new(commont)

type commont struct{}

// 哈希加密
//
// GET /hash
func (s *commont) Hash(c *gin.Context) {
	c.String(200, "commont List")
}
