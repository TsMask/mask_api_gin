package controller

import (
	"fmt"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/constants/token"
	"mask_api_gin/src/framework/model/result"
	"mask_api_gin/src/framework/service/ctx"
	tokenService "mask_api_gin/src/framework/service/token"
	commonModel "mask_api_gin/src/modules/common/model"
	commonService "mask_api_gin/src/modules/common/service"
	monitorService "mask_api_gin/src/modules/monitor/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 账号身份操作处理
var Account = &accountController{
	accountService:       commonService.AccountImpl,
	sysLogininforService: monitorService.SysLogininforImpl,
}

type accountController struct {
	// 账号身份操作服务
	accountService commonService.IAccount
	// 系统登录访问
	sysLogininforService monitorService.ISysLogininfor
}

// 系统登录
//
// POST /login
func (s *accountController) Login(c *gin.Context) {
	var loginBody commonModel.LoginBody
	if err := c.ShouldBindJSON(&loginBody); err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	// 必要字段
	if loginBody.Username == "" || loginBody.Password == "" {
		c.JSON(200, result.Err(nil))
		return
	}

	// 当前请求信息
	ipaddr, location := ctx.ClientIP(c)
	os, browser := ctx.UaOsBrowser(c)

	// 校验验证码
	err := s.accountService.ValidateCaptcha(
		loginBody.Username,
		loginBody.Code,
		loginBody.UUID,
	)
	// 根据错误信息，创建系统访问记录
	if err != nil {
		msg := err.Error() + " " + loginBody.Code
		s.sysLogininforService.NewLogininfor(
			loginBody.Username, common.STATUS_NO, msg,
			ipaddr, location, os, browser,
		)
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	// 登录用户信息
	loginUser, err := s.accountService.LoginByUsername(loginBody.Username, loginBody.Password)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	// 生成令牌，创建系统访问记录
	tokenStr := tokenService.Create(&loginUser, ipaddr, location, os, browser)
	if tokenStr == "" {
		c.JSON(200, result.Err(nil))
		return
	} else {
		s.sysLogininforService.NewLogininfor(
			loginBody.Username, common.STATUS_YES, "登录成功",
			ipaddr, location, os, browser,
		)
	}

	c.JSON(200, result.OkData(map[string]interface{}{
		token.RESPONSE_FIELD: tokenStr,
	}))
}

// 登录用户信息
//
// GET /getInfo
func (s *accountController) Info(c *gin.Context) {
	name := viper.GetString("pkg.name")
	version := viper.GetString("pkg.version")
	str := "欢迎使用%s后台管理框架，当前版本：%s，请通过前端管理地址访问。"
	c.JSON(200, result.OkMsg(fmt.Sprintf(str, name, version)))
}

// 登录用户路由信息
//
// GET /getRouters
func (s *accountController) Router(c *gin.Context) {
	name := viper.GetString("pkg.name")
	version := viper.GetString("pkg.version")
	str := "欢迎使用%s后台管理框架，当前版本：%s，请通过前端管理地址访问。"
	c.JSON(200, result.OkMsg(fmt.Sprintf(str, name, version)))
}

// 系统登出
//
// POST /logout
func (s *accountController) Logout(c *gin.Context) {
	c.JSON(200, result.OkMsg("退出成功"))
}
