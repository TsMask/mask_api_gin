package controller

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/gin-gonic/gin"
)

// NewCommon 实例化控制层
var NewCommon = &CommonController{}

// CommonController 通用请求 控制层处理
//
// PATH /
type CommonController struct{}

// Hash 哈希编码
//
// POST /hash
func (s *CommonController) Hash(c *gin.Context) {
	var body struct {
		Type string `json:"type" binding:"required,oneof=sha1 sha256 sha512 md5"`
		Str  string `json:"str" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var h hash.Hash
	var err error
	switch body.Type {
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	case "md5":
		h = md5.New()
	default:
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  fmt.Sprintf("不支持的哈希算法: %s", body.Type),
		})
		return
	}

	// 写入需要哈希的数据
	if _, err = h.Write([]byte(body.Str)); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "哈希写入错误",
		})
		return
	}

	// 计算哈希值的16进制表示
	hashed := h.Sum(nil)
	text := hex.EncodeToString(hashed)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": text,
	})
}
