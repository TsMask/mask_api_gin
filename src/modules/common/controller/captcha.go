package controller

import (
	"mask_api_gin/src/framework/model/result"

	"github.com/gin-gonic/gin"
)

// 验证码操作处理
var Captcha = new(captcha)

type captcha struct{}

// 获取验证码
//
// GET /captchaImage
func (s *captcha) CaptchaImage(c *gin.Context) {
	c.JSON(200, result.OkMsg("sdnfo"))
}
