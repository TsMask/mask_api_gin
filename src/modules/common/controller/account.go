package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	constSystem "mask_api_gin/src/framework/constants/system"
	constToken "mask_api_gin/src/framework/constants/token"
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	tokenUtils "mask_api_gin/src/framework/utils/token"
	"mask_api_gin/src/framework/vo/result"
	commonModel "mask_api_gin/src/modules/common/model"
	commonService "mask_api_gin/src/modules/common/service"
	systemService "mask_api_gin/src/modules/system/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewAccount 实例化控制层
var NewAccount = &AccountController{
	accountService:     commonService.NewAccount,
	sysLogLoginService: systemService.NewSysLogLogin,
}

// AccountController 账号身份操作 控制层处理
//
// PATH /
type AccountController struct {
	accountService     *commonService.Account            // 账号身份操作服务
	sysLogLoginService systemService.ISysLogLoginService // 系统登录访问
}

// Login 系统登录
//
// POST /login
func (s AccountController) Login(c *gin.Context) {
	var loginBody commonModel.LoginBody
	if err := c.ShouldBindJSON(&loginBody); err != nil {
		c.JSON(400, result.CodeMsg(40010, "参数错误"))
		return
	}

	// 当前请求信息
	ipaddr, location := ctxUtils.IPAddrLocation(c)
	os, browser := ctxUtils.UaOsBrowser(c)

	// 校验验证码 根据错误信息，创建系统访问记录
	if err := s.accountService.ValidateCaptcha(loginBody.Code, loginBody.UUID); err != nil {
		msg := fmt.Sprintf("%s code: %s", err.Error(), loginBody.Code)
		s.sysLogLoginService.Insert(
			loginBody.Username, constSystem.STATUS_NO, msg,
			[4]string{ipaddr, location, os, browser},
		)
		c.JSON(400, result.CodeMsg(40012, err.Error()))
		return
	}

	// 登录用户信息
	loginUser, err := s.accountService.ByUsername(loginBody.Username, loginBody.Password)
	if err != nil {
		c.JSON(400, result.CodeMsg(40014, err.Error()))
		return
	}

	// 生成令牌，创建系统访问记录
	tokenStr := tokenUtils.Create(&loginUser, [4]string{ipaddr, location, os, browser})
	if tokenStr == "" {
		c.JSON(400, result.CodeMsg(40001, "生成token失败"))
		return
	} else {
		s.accountService.UpdateLoginDateAndIP(&loginUser)
		s.sysLogLoginService.Insert(
			loginBody.Username, constSystem.STATUS_YES, "登录成功",
			[4]string{ipaddr, location, os, browser},
		)
	}

	c.JSON(200, result.OkData(map[string]any{
		"access_token": tokenStr,
		"token_type":   strings.TrimRight(constToken.HEADER_PREFIX, " "),
		"expires_in":   (loginUser.ExpireTime - loginUser.LoginTime) / 1000,
		"user_id":      loginUser.UserID,
	}))
}

// Me 登录用户信息
//
// GET /me
func (s AccountController) Me(c *gin.Context) {
	loginUser, err := ctxUtils.LoginUser(c)
	if err != nil {
		c.JSON(401, result.CodeMsg(401, err.Error()))
		return
	}

	// 角色权限集合，系统管理员拥有所有权限
	isSysAdmin := config.IsSysAdmin(loginUser.UserID)
	roles, perms := s.accountService.RoleAndMenuPerms(loginUser.UserID, isSysAdmin)

	c.JSON(200, result.OkData(map[string]any{
		"user":        loginUser.User,
		"roles":       roles,
		"permissions": perms,
	}))
}

// Router 登录用户路由信息
//
// GET /router
func (s AccountController) Router(c *gin.Context) {
	userID := ctxUtils.LoginUserToUserID(c)

	// 前端路由，系统管理员拥有所有
	isSysAdmin := config.IsSysAdmin(userID)
	buildMenus := s.accountService.RouteMenus(userID, isSysAdmin)
	c.JSON(200, result.OkData(buildMenus))
}

// Logout 系统登出
//
// POST /logout
func (s AccountController) Logout(c *gin.Context) {
	tokenStr := ctxUtils.Authorization(c)
	if tokenStr != "" {
		// 存在token时记录退出信息
		userName := tokenUtils.Remove(tokenStr)
		if userName != "" {
			// 当前请求信息
			ipaddr, location := ctxUtils.IPAddrLocation(c)
			os, browser := ctxUtils.UaOsBrowser(c)
			// 创建系统访问记录
			s.sysLogLoginService.Insert(
				userName, constSystem.STATUS_YES, "退出成功",
				[4]string{ipaddr, location, os, browser},
			)
		}
	}

	c.JSON(200, result.OkMsg("退出成功"))
}
