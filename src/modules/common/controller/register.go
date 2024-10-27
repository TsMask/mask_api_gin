package controller

import (
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/framework/vo/result"
	commonModel "mask_api_gin/src/modules/common/model"
	commonService "mask_api_gin/src/modules/common/service"
	systemService "mask_api_gin/src/modules/system/service"

	"github.com/gin-gonic/gin"
)

// NewRegister 实例化控制层
var NewRegister = &RegisterController{
	registerService:    commonService.NewRegisterService,
	sysLogLoginService: systemService.NewSysLogLogin,
}

// RegisterController 账号注册操作 控制层处理
//
// PATH /
type RegisterController struct {
	registerService    commonService.IRegisterService // 账号注册操作服务
	sysLogLoginService *systemService.SysLogLogin     // 系统登录访问服务
}

// Register 账号注册
//
// POST /register
func (s *RegisterController) Register(c *gin.Context) {
	var registerBody commonModel.RegisterBody
	if err := c.ShouldBindJSON(&registerBody); err != nil {
		c.JSON(400, result.ErrMsg("参数错误"))
		return
	}

	// 当前请求信息
	ipaddr, location := ctx.IPAddrLocation(c)
	os, browser := ctx.UaOsBrowser(c)

	// 校验验证码
	err := s.registerService.ValidateCaptcha(
		registerBody.Code,
		registerBody.UUID,
	)
	// 根据错误信息，创建系统访问记录
	if err != nil {
		msg := err.Error() + " " + registerBody.Code
		s.sysLogLoginService.Insert(
			registerBody.Username, constSystem.STATUS_NO, msg,
			[4]string{ipaddr, location, os, browser},
		)
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	// 判断必传参数
	if !regular.ValidUsername(registerBody.Username) {
		c.JSON(200, result.ErrMsg("用户账号只能包含大写小写字母，数字，且不少于6位"))
		return
	}
	if !regular.ValidPassword(registerBody.Password) {
		c.JSON(200, result.ErrMsg("登录密码至少包含大小写字母、数字、特殊符号，且不少于6位"))
		return
	}
	if registerBody.Password != registerBody.ConfirmPassword {
		c.JSON(200, result.ErrMsg("用户确认输入密码不一致"))
		return
	}

	userID, err := s.registerService.ByUserName(registerBody.Username, registerBody.Password, registerBody.UserType)
	if err == nil {
		msg := registerBody.Username + " 注册成功 " + userID
		s.sysLogLoginService.Insert(
			registerBody.Username, constSystem.STATUS_YES, msg,
			[4]string{ipaddr, location, os, browser},
		)
		c.JSON(200, result.OkMsg("注册成功"))
		return
	}
	c.JSON(200, result.ErrMsg(err.Error()))
}
