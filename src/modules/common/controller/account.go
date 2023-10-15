package controller

import (
	"mask_api_gin/src/framework/config"
	commonConstants "mask_api_gin/src/framework/constants/common"
	tokenConstants "mask_api_gin/src/framework/constants/token"
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	tokenUtils "mask_api_gin/src/framework/utils/token"
	"mask_api_gin/src/framework/vo/result"
	commonModel "mask_api_gin/src/modules/common/model"
	commonService "mask_api_gin/src/modules/common/service"
	systemService "mask_api_gin/src/modules/system/service"

	"github.com/gin-gonic/gin"
)

// 实例化控制层 AccountController 结构体
var NewAccount = &AccountController{
	accountService:     commonService.NewAccountImpl,
	sysLogLoginService: systemService.NewSysLogLoginImpl,
}

// 账号身份操作处理
//
// PATH /
type AccountController struct {
	// 账号身份操作服务
	accountService commonService.IAccount
	// 系统登录访问
	sysLogLoginService systemService.ISysLogLogin
}

// 系统登录
//
// POST /login
func (s *AccountController) Login(c *gin.Context) {
	var loginBody commonModel.LoginBody
	if err := c.ShouldBindJSON(&loginBody); err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 当前请求信息
	ipaddr, location := ctxUtils.IPAddrLocation(c)
	os, browser := ctxUtils.UaOsBrowser(c)

	// 校验验证码
	err := s.accountService.ValidateCaptcha(
		loginBody.Code,
		loginBody.UUID,
	)
	// 根据错误信息，创建系统访问记录
	if err != nil {
		msg := err.Error() + " " + loginBody.Code
		s.sysLogLoginService.NewSysLogLogin(
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
		s.sysLogLoginService.NewSysLogLogin(
			loginBody.Username, commonConstants.STATUS_YES, "登录成功",
			ipaddr, location, os, browser,
		)
	}

	c.JSON(200, result.OkData(map[string]any{
		tokenConstants.RESPONSE_FIELD: tokenStr,
	}))
}

// 登录用户信息
//
// GET /getInfo
func (s *AccountController) Info(c *gin.Context) {
	loginUser, err := ctxUtils.LoginUser(c)
	if err != nil {
		c.JSON(401, result.CodeMsg(401, err.Error()))
		return
	}

	// 角色权限集合，管理员拥有所有权限
	isAdmin := config.IsAdmin(loginUser.UserID)
	roles, perms := s.accountService.RoleAndMenuPerms(loginUser.UserID, isAdmin)

	c.JSON(200, result.OkData(map[string]any{
		"user":        loginUser.User,
		"roles":       roles,
		"permissions": perms,
	}))
}

// 登录用户路由信息
//
// GET /getRouters
func (s *AccountController) Router(c *gin.Context) {
	userID := ctxUtils.LoginUserToUserID(c)

	// 前端路由，管理员拥有所有
	isAdmin := config.IsAdmin(userID)
	buildMenus := s.accountService.RouteMenus(userID, isAdmin)
	c.JSON(200, result.OkData(buildMenus))
}

// 系统登出
//
// POST /logout
func (s *AccountController) Logout(c *gin.Context) {
	tokenStr := ctxUtils.Authorization(c)
	if tokenStr != "" {
		// 存在token时记录退出信息
		userName := tokenUtils.Remove(tokenStr)
		if userName != "" {
			// 当前请求信息
			ipaddr, location := ctxUtils.IPAddrLocation(c)
			os, browser := ctxUtils.UaOsBrowser(c)
			// 创建系统访问记录
			s.sysLogLoginService.NewSysLogLogin(
				userName, commonConstants.STATUS_NO, "退出成功",
				ipaddr, location, os, browser,
			)
		}
	}

	c.JSON(200, result.OkMsg("退出成功"))
}
