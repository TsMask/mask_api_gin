package controller

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/token"
	"mask_api_gin/src/modules/common/model"
	"mask_api_gin/src/modules/common/service"
	systemService "mask_api_gin/src/modules/system/service"
)

// NewAccount 实例化控制层
var NewAccount = &AccountController{
	accountService:     service.NewAccount,
	sysLogLoginService: systemService.NewSysLogLogin,
}

// AccountController 账号身份操作 控制层处理
//
// PATH /
type AccountController struct {
	accountService     *service.Account           // 账号身份操作服务
	sysLogLoginService *systemService.SysLogLogin // 系统登录访问
}

// Login 系统登录
//
// POST /login
func (s AccountController) Login(c *gin.Context) {
	var body model.LoginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(40422, errMsgs))
		return
	}

	// 当前请求信息
	ipaddr, location := reqctx.IPAddrLocation(c)
	os, browser := reqctx.UaOsBrowser(c)

	// 校验验证码 根据错误信息，创建系统访问记录
	if err := s.accountService.ValidateCaptcha(body.Code, body.UUID); err != nil {
		msg := fmt.Sprintf("%s code: %s", err.Error(), body.Code)
		s.sysLogLoginService.Insert(
			body.Username, constants.STATUS_NO, msg,
			[4]string{ipaddr, location, os, browser},
		)
		c.JSON(400, resp.CodeMsg(40012, err.Error()))
		return
	}

	// 登录用户信息
	info, err := s.accountService.ByUsername(body.Username, body.Password)
	if err != nil {
		c.JSON(400, resp.CodeMsg(40014, err.Error()))
		return
	}

	// 生成令牌，创建系统访问记录
	tokenStr := token.Create(&info, [4]string{ipaddr, location, os, browser})
	if tokenStr == "" {
		c.JSON(400, resp.CodeMsg(40001, "token generation failed"))
		return
	} else {
		s.accountService.UpdateLoginDateAndIP(info)
		s.sysLogLoginService.Insert(
			body.Username, constants.STATUS_YES, "登录成功",
			[4]string{ipaddr, location, os, browser},
		)
	}

	c.JSON(200, resp.OkData(map[string]any{
		"accessToken": tokenStr,
		"tokenType":   strings.TrimRight(constants.HEADER_PREFIX, " "),
		"expiresIn":   (info.ExpireTime - info.LoginTime) / 1000,
		"userId":      info.UserId,
	}))
}

// Me 登录用户信息
//
// GET /me
func (s AccountController) Me(c *gin.Context) {
	info, err := reqctx.LoginUser(c)
	if err != nil {
		c.JSON(401, resp.CodeMsg(40003, err.Error()))
		return
	}

	// 角色权限集合，系统管理员拥有所有权限
	isSystemUser := config.IsSystemUser(info.UserId)
	roles, perms := s.accountService.RoleAndMenuPerms(info.UserId, isSystemUser)

	c.JSON(200, resp.OkData(map[string]any{
		"user":        info.User,
		"roles":       roles,
		"permissions": perms,
	}))
}

// Router 登录用户路由信息
//
// GET /router
func (s AccountController) Router(c *gin.Context) {
	loginUserId := reqctx.LoginUserToUserID(c)

	// 前端路由，系统管理员拥有所有
	isSystemUser := config.IsSystemUser(loginUserId)
	buildMenus := s.accountService.RouteMenus(loginUserId, isSystemUser)
	c.JSON(200, resp.OkData(buildMenus))
}

// Logout 系统登出
//
// POST /logout
func (s AccountController) Logout(c *gin.Context) {
	tokenStr := reqctx.Authorization(c)
	if tokenStr != "" {
		// 存在token时记录退出信息
		userName := token.Remove(tokenStr)
		if userName != "" {
			// 当前请求信息
			ipaddr, location := reqctx.IPAddrLocation(c)
			os, browser := reqctx.UaOsBrowser(c)
			// 创建系统访问记录
			s.sysLogLoginService.Insert(
				userName, constants.STATUS_YES, "退出成功",
				[4]string{ipaddr, location, os, browser},
			)
		}
	}
	c.JSON(200, resp.OkMsg("logout successful"))
}
