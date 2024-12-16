package controller

import (
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/context"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/token"
	commonModel "mask_api_gin/src/modules/common/model"
	commonService "mask_api_gin/src/modules/common/service"
	systemService "mask_api_gin/src/modules/system/service"

	"fmt"
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
	accountService     *commonService.Account     // 账号身份操作服务
	sysLogLoginService *systemService.SysLogLogin // 系统登录访问
}

// Login 系统登录
//
// POST /login
func (s AccountController) Login(c *gin.Context) {
	var loginBody commonModel.LoginBody
	if err := c.ShouldBindJSON(&loginBody); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}

	// 当前请求信息
	ipaddr, location := context.IPAddrLocation(c)
	os, browser := context.UaOsBrowser(c)

	// 校验验证码 根据错误信息，创建系统访问记录
	if err := s.accountService.ValidateCaptcha(loginBody.Code, loginBody.UUID); err != nil {
		msg := fmt.Sprintf("%s code: %s", err.Error(), loginBody.Code)
		s.sysLogLoginService.Insert(
			loginBody.Username, constants.STATUS_NO, msg,
			[4]string{ipaddr, location, os, browser},
		)
		c.JSON(400, response.CodeMsg(40012, err.Error()))
		return
	}

	// 登录用户信息
	loginUser, err := s.accountService.ByUsername(loginBody.Username, loginBody.Password)
	if err != nil {
		c.JSON(400, response.CodeMsg(40014, err.Error()))
		return
	}

	// 生成令牌，创建系统访问记录
	tokenStr := token.Create(&loginUser, [4]string{ipaddr, location, os, browser})
	if tokenStr == "" {
		c.JSON(400, response.CodeMsg(40001, "token generation failed"))
		return
	} else {
		s.accountService.UpdateLoginDateAndIP(&loginUser)
		s.sysLogLoginService.Insert(
			loginBody.Username, constants.STATUS_YES, "登录成功",
			[4]string{ipaddr, location, os, browser},
		)
	}

	c.JSON(200, response.OkData(map[string]any{
		"accessToken": tokenStr,
		"tokenType":   strings.TrimRight(constants.HEADER_PREFIX, " "),
		"expiresIn":   (loginUser.ExpireTime - loginUser.LoginTime) / 1000,
		"userId":      loginUser.UserId,
	}))
}

// Me 登录用户信息
//
// GET /me
func (s AccountController) Me(c *gin.Context) {
	loginUser, err := context.LoginUser(c)
	if err != nil {
		c.JSON(401, response.CodeMsg(40003, err.Error()))
		return
	}

	// 角色权限集合，系统管理员拥有所有权限
	isSystemUser := config.IsSystemUser(loginUser.UserId)
	roles, perms := s.accountService.RoleAndMenuPerms(loginUser.UserId, isSystemUser)

	c.JSON(200, response.OkData(map[string]any{
		"user":        loginUser.User,
		"roles":       roles,
		"permissions": perms,
	}))
}

// Router 登录用户路由信息
//
// GET /router
func (s AccountController) Router(c *gin.Context) {
	userId := context.LoginUserToUserID(c)

	// 前端路由，系统管理员拥有所有
	isSystemUser := config.IsSystemUser(userId)
	buildMenus := s.accountService.RouteMenus(userId, isSystemUser)
	c.JSON(200, response.OkData(buildMenus))
}

// Logout 系统登出
//
// POST /logout
func (s AccountController) Logout(c *gin.Context) {
	tokenStr := context.Authorization(c)
	if tokenStr != "" {
		// 存在token时记录退出信息
		userName := token.Remove(tokenStr)
		if userName != "" {
			// 当前请求信息
			ipaddr, location := context.IPAddrLocation(c)
			os, browser := context.UaOsBrowser(c)
			// 创建系统访问记录
			s.sysLogLoginService.Insert(
				userName, constants.STATUS_YES, "退出成功",
				[4]string{ipaddr, location, os, browser},
			)
		}
	}
	c.JSON(200, response.OkMsg("logout successful"))
}
