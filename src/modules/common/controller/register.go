package controller

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/regular"
	commonModel "mask_api_gin/src/modules/common/model"
	commonService "mask_api_gin/src/modules/common/service"
	systemService "mask_api_gin/src/modules/system/service"

	"fmt"

	"github.com/gin-gonic/gin"
)

// NewRegister 实例化控制层
var NewRegister = &RegisterController{
	registerService:    commonService.NewRegister,
	sysLogLoginService: systemService.NewSysLogLogin,
}

// RegisterController 账号注册操作 控制层处理
//
// PATH /
type RegisterController struct {
	registerService    *commonService.Register    // 账号注册操作服务
	sysLogLoginService *systemService.SysLogLogin // 系统登录访问服务
}

// Register 账号注册
//
// POST /register
func (s *RegisterController) Register(c *gin.Context) {
	var body commonModel.RegisterBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}

	// 当前请求信息
	ipaddr, location := ctx.IPAddrLocation(c)
	os, browser := ctx.UaOsBrowser(c)

	// 校验验证码
	err := s.registerService.ValidateCaptcha(body.Code, body.UUID)
	// 根据错误信息，创建系统访问记录
	if err != nil {
		msg := err.Error() + " " + body.Code
		s.sysLogLoginService.Insert(
			body.Username, constants.STATUS_NO, msg,
			[4]string{ipaddr, location, os, browser},
		)
		c.JSON(400, response.CodeMsg(40012, err.Error()))
		return
	}

	// 判断必传参数
	if !regular.ValidUsername(body.Username) {
		c.JSON(200, response.ErrMsg("用户账号只能包含大写小写字母，数字，且不少于4位"))
		return
	}
	if !regular.ValidPassword(body.Password) {
		c.JSON(200, response.ErrMsg("登录密码至少包含大小写字母、数字、特殊符号，且不少于6位"))
		return
	}
	if body.Password != body.ConfirmPassword {
		c.JSON(200, response.ErrMsg("用户确认输入密码不一致"))
		return
	}

	userId, err := s.registerService.ByUserName(body.Username, body.Password)
	if err == nil {
		msg := fmt.Sprintf("%s 注册成功 %s", body.Username, userId)
		s.sysLogLoginService.Insert(
			body.Username, constants.STATUS_YES, msg,
			[4]string{ipaddr, location, os, browser},
		)
		c.JSON(200, response.OkMsg("注册成功"))
		return
	}
	c.JSON(200, response.ErrMsg(err.Error()))
}
