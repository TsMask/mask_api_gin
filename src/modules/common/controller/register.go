package controller

import (
	commonConstants "mask_api_gin/src/framework/constants/common"
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/framework/vo/result"
	commonModel "mask_api_gin/src/modules/common/model"
	commonService "mask_api_gin/src/modules/common/service"
	systemService "mask_api_gin/src/modules/system/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// 实例化控制层 RegisterController 结构体
var NewRegister = &RegisterController{
	registerService:    commonService.NewRegisterImpl,
	sysLogLoginService: systemService.NewSysLogLoginImpl,
}

// 账号注册操作处理
//
// PATH /
type RegisterController struct {
	// 账号注册操作服务
	registerService commonService.IRegister
	// 系统登录访问
	sysLogLoginService systemService.ISysLogLogin
}

// 账号注册
//
// GET /captchaImage
func (s *RegisterController) UserName(c *gin.Context) {
	var registerBody commonModel.RegisterBody
	if err := c.ShouldBindJSON(&registerBody); err != nil {
		c.JSON(400, result.ErrMsg("参数错误"))
		return
	}

	// 判断必传参数
	if !regular.ValidUsername(registerBody.Username) {
		c.JSON(200, result.ErrMsg("账号不能以数字开头，可包含大写小写字母，数字，且不少于5位"))
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

	// 当前请求信息
	ipaddr, location := ctxUtils.IPAddrLocation(c)
	os, browser := ctxUtils.UaOsBrowser(c)

	// 校验验证码
	err := s.registerService.ValidateCaptcha(
		registerBody.Code,
		registerBody.UUID,
	)
	// 根据错误信息，创建系统访问记录
	if err != nil {
		msg := err.Error() + " " + registerBody.Code
		s.sysLogLoginService.CreateSysLogLogin(
			registerBody.Username, commonConstants.STATUS_NO, msg,
			ipaddr, location, os, browser,
		)
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	infoStr := s.registerService.ByUserName(registerBody.Username, registerBody.Password, registerBody.UserType)
	if !strings.HasPrefix(infoStr, "注册") {
		msg := registerBody.Username + " 注册成功 " + infoStr
		s.sysLogLoginService.CreateSysLogLogin(
			registerBody.Username, commonConstants.STATUS_YES, msg,
			ipaddr, location, os, browser,
		)
		c.JSON(200, result.OkMsg("注册成功"))
		return
	}
	c.JSON(200, result.ErrMsg(infoStr))
}
