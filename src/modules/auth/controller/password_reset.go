package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/regular"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/modules/auth/model"
	"mask_api_gin/src/modules/auth/service"
)

// NewPasswordReset 实例化控制层
var NewPasswordReset = &PasswordResetController{
	passwordResetService: service.NewPasswordReset,
}

// PasswordResetController 找回密码操作 控制层处理
//
// PATH /api/v1/auth/password-reset
type PasswordResetController struct {
	passwordResetService *service.PasswordReset // 找回密码服务
}

// Lookup 查询密码找回联系方式
//
// POST /lookup
func (s PasswordResetController) Lookup(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"` // 用户账号
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsg := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsg))
		return
	}

	contact, err := s.passwordResetService.Lookup(body.Username)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, resp.OkData(contact))
}

// CaptchaEmail 发送密码找回验证码
//
// POST /captcha-email
func (s PasswordResetController) CaptchaEmail(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"` // 用户账号
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsg := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsg))
		return
	}

	captcha, err := s.passwordResetService.SendCaptchaEmail(body.Username)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, resp.OkData(captcha))
}

// Reset 重置登录密码
//
// POST /
func (s PasswordResetController) Reset(c *gin.Context) {
	var body model.PasswordResetBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsg := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsg))
		return
	}

	if !regular.ValidPassword(body.Password) {
		c.JSON(200, resp.ErrMsg("登录密码至少包含大小写字母、数字、特殊符号，且不少于6位"))
		return
	}
	if body.Password != body.ConfirmPassword {
		c.JSON(200, resp.ErrMsg("两次输入的新密码不一致"))
		return
	}

	if err := s.passwordResetService.Reset(body.Username, body.Code, body.UUID, body.Password); err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, resp.OkMsg("重置密码成功，请使用新密码登录"))
}
