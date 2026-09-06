package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/regular"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/modules/auth/model"
	"mask_api_gin/src/modules/auth/service"
	systemService "mask_api_gin/src/modules/system/service"
)

// NewRegister 实例化控制层
var NewRegister = &RegisterController{
	registerService:    service.NewRegister,
	sysLogLoginService: systemService.NewSysLogLogin,
}

// RegisterController 账号注册操作 控制层处理
//
// PATH /api/v1/auth
type RegisterController struct {
	registerService    *service.Register          // 账号注册操作服务
	sysLogLoginService *systemService.SysLogLogin // 系统登录访问服务
}

// Register 账号注册
//
// POST /register
func (s RegisterController) Register(c *gin.Context) {
	var body model.RegisterBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	if !regular.ValidUsername(body.Username) || regular.ValidMobile(body.Username) {
		c.JSON(200, resp.ErrMsg("用户账号格式不正确，请输入用户名格式"))
		return
	}
	if !regular.ValidPassword(body.Password) {
		c.JSON(200, resp.ErrMsg("登录密码格式不正确"))
		return
	}
	if body.Password != body.ConfirmPassword {
		c.JSON(200, resp.ErrMsg("确认密码不一致"))
		return
	}

	ipaddr, location := reqctx.IPAddrLocation(c)
	os, browser := reqctx.UaOsBrowser(c)
	ilobArr := [4]string{ipaddr, location, os, browser}
	if err := s.registerService.ValidateCaptcha(body.Code, body.UUID); err != nil {
		msg := fmt.Sprintf("%s code %s", err.Error(), body.Code)
		s.sysLogLoginService.Insert(body.Username, constants.STATUS_NO, msg, ilobArr)
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	userId, err := s.registerService.ByUserName(body.Username, body.Password)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	msg := fmt.Sprintf("%s 注册成功 %d", body.Username, userId)
	s.sysLogLoginService.Insert(body.Username, constants.STATUS_YES, msg, ilobArr)
	c.JSON(200, resp.OkMsg("注册成功"))
}

// RegisterPhone 手机号注册
//
// POST /register-phone
func (s RegisterController) RegisterPhone(c *gin.Context) {
	var body model.RegisterPhoneBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	if !regular.ValidMobile(body.Phone) {
		c.JSON(200, resp.ErrMsg("手机号注册时 phone 必须为有效手机号码"))
		return
	}

	ipaddr, location := reqctx.IPAddrLocation(c)
	os, browser := reqctx.UaOsBrowser(c)
	ilobArr := [4]string{ipaddr, location, os, browser}
	if err := s.registerService.ValidateSmsCaptcha(body.Phone, body.Code, body.UUID); err != nil {
		msg := fmt.Sprintf("%s code %s", err.Error(), body.Code)
		s.sysLogLoginService.Insert(body.Phone, constants.STATUS_NO, msg, ilobArr)
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	userId, err := s.registerService.ByPhone(body.Phone)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	msg := fmt.Sprintf("%s 手机号注册成功 %d", body.Phone, userId)
	s.sysLogLoginService.Insert(body.Phone, constants.STATUS_YES, msg, ilobArr)
	c.JSON(200, resp.OkMsg("注册成功"))
}
