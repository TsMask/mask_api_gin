package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/token"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/auth/model"
	"mask_api_gin/src/modules/auth/service"
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
// POST /auth/login
func (s AccountController) Login(c *gin.Context) {
	var body model.LoginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(422001, errMsgs))
		return
	}

	// 当前请求信息
	ipaddr, location := reqctx.IPAddrLocation(c)
	os, browser := reqctx.UaOsBrowser(c)

	// 校验验证码 根据错误信息，创建系统访问记录
	if err := s.accountService.ValidateCaptcha(body.Code, body.UUID); err != nil {
		msg := fmt.Sprintf("%s code %s", err.Error(), body.Code)
		s.sysLogLoginService.Insert(
			body.Username, constants.STATUS_NO, msg,
			[4]string{ipaddr, location, os, browser},
		)
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	// 登录用户信息
	info, err := s.accountService.ByUsername(body.Username, body.Password)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	deviceFingerprint := reqctx.DeviceFingerprint(c, info.UserId)

	// 生成访问令牌
	accessToken, expiresIn := token.UserTokenCreate(info.UserId, deviceFingerprint, "access")
	if accessToken == "" || expiresIn == 0 {
		c.JSON(200, resp.ErrMsg("token generation failed"))
		return
	}
	// 生成刷新令牌
	refreshToken, refreshExpiresIn := token.UserTokenCreate(info.UserId, deviceFingerprint, "refresh")

	// 记录令牌，创建系统访问记录
	token.UserInfoCreate(&info, deviceFingerprint, [4]string{ipaddr, location, os, browser})
	s.accountService.UpdateLoginDateAndIP(info)
	s.sysLogLoginService.Insert(
		body.Username, constants.STATUS_YES, "登录成功",
		[4]string{ipaddr, location, os, browser},
	)

	c.JSON(200, resp.OkData(map[string]any{
		"tokenType":        constants.HEADER_PREFIX,
		"accessToken":      accessToken,
		"expiresIn":        expiresIn,
		"refreshToken":     refreshToken,
		"refreshExpiresIn": refreshExpiresIn,
		"userId":           info.UserId,
	}))
}

// Logout 系统登出
//
// POST /auth/logout
func (s AccountController) Logout(c *gin.Context) {
	tokenStr := reqctx.Authorization(c)
	if tokenStr != "" {
		// 存在token时记录退出信息
		userName, err := token.UserInfoRemove(tokenStr)
		if err != nil {
			// 当前请求信息
			ipaddr, location := reqctx.IPAddrLocation(c)
			os, browser := reqctx.UaOsBrowser(c)
			// 创建系统访问记录
			s.sysLogLoginService.Insert(
				userName, constants.STATUS_YES, "主动注销登录信息",
				[4]string{ipaddr, location, os, browser},
			)
		}
	}
	c.JSON(200, resp.OkMsg("logout successful"))
}

// RefreshToken 刷新Token
//
// POST /auth/refresh-token
func (s AccountController) RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken" binding:"required"` // 刷新令牌
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(422001, errMsgs))
		return
	}

	// 验证刷新令牌是否有效
	claims, err := token.UserTokenVerify(body.RefreshToken, "refresh")
	if err != nil {
		c.JSON(401, resp.CodeMsg(401001, err.Error()))
		return
	}
	userId := parse.Number(claims[constants.JWT_USER_ID])

	// 登录用户信息
	info, err := s.accountService.ByUserId(userId)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	// 设备指纹信息是否一致
	deviceId := fmt.Sprint(claims[constants.JWT_DEVICE_ID])
	deviceFingerprint := reqctx.DeviceFingerprint(c, userId)
	if deviceId != deviceFingerprint {
		c.JSON(200, resp.ErrMsg("device fingerprint mismatch"))
		return
	}

	// 生成访问令牌
	accessToken, expiresIn := token.UserTokenCreate(userId, deviceFingerprint, "access")
	if accessToken == "" || expiresIn == 0 {
		c.JSON(200, resp.ErrMsg("token generation failed"))
		return
	}
	// 生成刷新令牌
	now := time.Now()
	exp, _ := claims.GetExpirationTime()
	iat, _ := claims.GetIssuedAt()
	refreshExpiresIn := int64(exp.Sub(now).Seconds())
	refreshToken := body.RefreshToken

	// 如果当前时间大于过期时间的一半，则生成新令牌
	halfExp := exp.Add(-(exp.Sub(iat.Time)) / 2)
	if now.After(halfExp) {
		refreshToken, refreshExpiresIn = token.UserTokenCreate(userId, deviceFingerprint, "refresh")
	}

	// 当前请求信息
	ipaddr, location := reqctx.IPAddrLocation(c)
	os, browser := reqctx.UaOsBrowser(c)
	// 记录令牌，创建系统访问记录
	token.UserInfoCreate(&info, deviceFingerprint, [4]string{ipaddr, location, os, browser})
	s.accountService.UpdateLoginDateAndIP(info)
	s.sysLogLoginService.Insert(
		info.User.UserName, constants.STATUS_YES, "刷新访问令牌成功",
		[4]string{ipaddr, location, os, browser},
	)

	// 返回访问令牌和刷新令牌
	c.JSON(200, resp.OkData(map[string]any{
		"tokenType":        constants.HEADER_PREFIX,
		"accessToken":      accessToken,
		"expiresIn":        expiresIn,
		"refreshToken":     refreshToken,
		"refreshExpiresIn": refreshExpiresIn,
		"userId":           userId,
	}))
}

// Me 登录用户信息
//
// GET /me
func (s AccountController) Me(c *gin.Context) {
	info, err := reqctx.LoginUser(c)
	if err != nil {
		c.JSON(401, resp.CodeMsg(401002, err.Error()))
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
