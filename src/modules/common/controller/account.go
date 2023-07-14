package controller

import (
	"mask_api_gin/src/framework/config"
	commonConstants "mask_api_gin/src/framework/constants/common"
	tokenConstants "mask_api_gin/src/framework/constants/token"
	"mask_api_gin/src/framework/model/result"
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	tokenUtils "mask_api_gin/src/framework/utils/token"
	commonModel "mask_api_gin/src/modules/common/model"
	commonService "mask_api_gin/src/modules/common/service"
	monitorService "mask_api_gin/src/modules/monitor/service"

	"github.com/gin-gonic/gin"
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
	ipaddr, location := ctxUtils.IPAddrLocation(c)
	os, browser := ctxUtils.UaOsBrowser(c)

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
			loginBody.Username, commonConstants.STATUS_NO, msg,
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
	tokenStr := tokenUtils.Create(&loginUser, ipaddr, location, os, browser)
	if tokenStr == "" {
		c.JSON(200, result.Err(nil))
		return
	} else {
		s.sysLogininforService.NewLogininfor(
			loginBody.Username, commonConstants.STATUS_YES, "登录成功",
			ipaddr, location, os, browser,
		)
	}

	c.JSON(200, result.OkData(map[string]interface{}{
		tokenConstants.RESPONSE_FIELD: tokenStr,
	}))
}

// 登录用户信息
//
// GET /getInfo
func (s *accountController) Info(c *gin.Context) {
	loginUser, err := ctxUtils.LoginUser(c)
	if err != nil {
		c.JSON(401, result.ErrMsg(err.Error()))
		return
	}

	// 角色权限集合，管理员拥有所有权限
	isAdmin := config.IsAdmin(loginUser.UserID)
	roles, perms := s.accountService.RoleAndMenuPerms(loginUser.UserID, isAdmin)

	c.JSON(200, result.OkData(map[string]interface{}{
		"user":        loginUser.User,
		"roles":       roles,
		"permissions": perms,
	}))
}

// 登录用户路由信息
//
// GET /getRouters
func (s *accountController) Router(c *gin.Context) {
	loginUser, err := ctxUtils.LoginUser(c)
	if err != nil {
		c.JSON(401, result.ErrMsg(err.Error()))
		return
	}

	// 前端路由，管理员拥有所有
	isAdmin := config.IsAdmin(loginUser.UserID)
	buildMenus := s.accountService.RouteMenus(loginUser.UserID, isAdmin)
	c.JSON(200, result.OkData(buildMenus))
}

// 系统登出
//
// POST /logout
func (s *accountController) Logout(c *gin.Context) {
	tokenStr := ctxUtils.Authorization(c)
	if tokenStr != "" {
		// 存在token时记录退出信息
		userName := tokenUtils.Remove(tokenStr)
		if userName != "" {
			// 当前请求信息
			ipaddr, location := ctxUtils.IPAddrLocation(c)
			os, browser := ctxUtils.UaOsBrowser(c)
			// 创建系统访问记录
			s.sysLogininforService.NewLogininfor(
				userName, commonConstants.STATUS_NO, "退出成功",
				ipaddr, location, os, browser,
			)
		}
	}

	c.JSON(200, result.OkMsg("退出成功"))
}
