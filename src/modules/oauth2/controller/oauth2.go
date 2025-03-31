package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/token"
	"mask_api_gin/src/modules/oauth2/model"
	"mask_api_gin/src/modules/oauth2/service"
)

// NewOauth2 实例化控制层
var NewOauth2 = &Oauth2Controller{
	oauth2Service:         service.NewOauth2Service,
	oauth2ClientService:   service.NewOauth2ClientService,
	oauth2LogLoginService: service.NewOauth2LogLogin,
}

// Oauth2Controller 授权第三方客户端应用认证 控制层处理
//
// PATH /oauth2
type Oauth2Controller struct {
	oauth2Service         *service.Oauth2Service         // 用户授权第三方信息服务
	oauth2ClientService   *service.Oauth2ClientService   // 用户授权第三方应用信息服务
	oauth2LogLoginService *service.Oauth2LogLoginService // 用户授权第三方应用登录日志
}

// Authorize 获取登录预授权码
//
// GET /authorize
func (s Oauth2Controller) Authorize(c *gin.Context) {
	var query model.CodeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(40422, errMsgs))
		return
	}

	// 是否存在clientId
	info := s.oauth2ClientService.FindByClientId(query.ClientId)
	if info.ClientId == "" || info.ClientId != query.ClientId {
		c.JSON(400, resp.CodeMsg(40404, "clientId不存在"))
		return
	}

	// 判断IP白名单
	if !strings.Contains(info.IPWhite, c.ClientIP()) {
		c.JSON(400, resp.CodeMsg(40404, "IP白名单不匹配"))
		return
	}

	// 生成登录预授权码
	code := s.oauth2Service.CreateCode()

	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", query.RedirectUrl, code, query.State)
	c.Redirect(302, redirectURL)
}

// Tooken 通过授权码获取访问令牌
//
// POST /token
func (s Oauth2Controller) Token(c *gin.Context) {
	var body model.TokenBody
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(40422, errMsgs))
		return
	}
	if body.GrantType != "authorization_code" || body.Code == "" {
		c.JSON(400, resp.CodeMsg(40404, "grantType or code error"))
		return
	}

	// 当前请求信息
	ipaddr, location := reqctx.IPAddrLocation(c)
	os, browser := reqctx.UaOsBrowser(c)

	// 校验验证码 根据错误信息，创建系统访问记录
	if err := s.oauth2Service.ValidateCode(body.Code); err != nil {
		msg := fmt.Sprintf("%s code %s", err.Error(), body.Code)
		s.oauth2LogLoginService.Insert(
			body.ClientId, constants.STATUS_NO, msg,
			[4]string{ipaddr, location, os, browser},
		)
		c.JSON(400, resp.CodeMsg(40012, err.Error()))
		return
	}

	// 客户端信息
	info, err := s.oauth2Service.ByClient(body.ClientId, body.ClientSecret)
	if err != nil {
		c.JSON(400, resp.CodeMsg(40014, err.Error()))
		return
	}
	deviceFingerprint := reqctx.DeviceFingerprint(c, info.ClientId)

	// 生成访问令牌
	accessToken, expiresIn := token.Oauth2TokenCreate(info.ClientId, deviceFingerprint, "access")
	if accessToken == "" || expiresIn == 0 {
		c.JSON(400, resp.CodeMsg(40001, "token generation failed"))
		return
	}
	// 生成刷新令牌
	refreshToken, refreshExpiresIn := token.Oauth2TokenCreate(info.ClientId, deviceFingerprint, "refresh")

	// 记录令牌，创建系统访问记录
	token.Oauth2InfoCreate(&info, deviceFingerprint, [4]string{ipaddr, location, os, browser})
	s.oauth2Service.UpdateLoginDateAndIP(info)
	s.oauth2LogLoginService.Insert(
		body.ClientId, constants.STATUS_YES, "授权成功",
		[4]string{ipaddr, location, os, browser},
	)

	c.JSON(200, resp.OkData(map[string]any{
		"tokenType":        constants.HEADER_PREFIX,
		"accessToken":      accessToken,
		"expiresIn":        expiresIn,
		"refreshToken":     refreshToken,
		"refreshExpiresIn": refreshExpiresIn,
	}))
}

// RefreshToken 通过刷新令牌续期访问令牌
//
// POST /refresh-token
func (s Oauth2Controller) RefreshToken(c *gin.Context) {
	var body model.TokenBody
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(40422, errMsgs))
		return
	}
	if body.GrantType != "refresh_token" || body.RefreshToken == "" {
		c.JSON(400, resp.CodeMsg(40404, "grantType or refreshToken error"))
		return
	}

	// 验证刷新令牌是否有效
	claims, err := token.Oauth2TokenVerify(body.RefreshToken, "refresh")
	if err != nil {
		c.JSON(401, resp.CodeMsg(40001, err.Error()))
		return
	}
	clientId := fmt.Sprint(claims[constants.JWT_CLIENT_ID])

	// 客户端信息
	info, err := s.oauth2Service.ByClient(body.ClientId, body.ClientSecret)
	if err != nil {
		c.JSON(400, resp.CodeMsg(40014, err.Error()))
		return
	}

	// 客户端ID是否一致
	if clientId != body.ClientId {
		c.JSON(400, resp.CodeMsg(40003, "clientId mismatch"))
		return
	}
	// 设备指纹信息是否一致
	deviceId := fmt.Sprint(claims[constants.JWT_DEVICE_ID])
	deviceFingerprint := reqctx.DeviceFingerprint(c, clientId)
	if deviceId != deviceFingerprint {
		c.JSON(401, resp.CodeMsg(40003, "device fingerprint mismatch"))
		return
	}

	// 生成访问令牌
	accessToken, expiresIn := token.Oauth2TokenCreate(clientId, deviceFingerprint, "access")
	if accessToken == "" || expiresIn == 0 {
		c.JSON(400, resp.CodeMsg(40001, "token generation failed"))
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
		refreshToken, refreshExpiresIn = token.Oauth2TokenCreate(clientId, deviceFingerprint, "refresh")
	}

	// 当前请求信息
	ipaddr, location := reqctx.IPAddrLocation(c)
	os, browser := reqctx.UaOsBrowser(c)
	// 记录令牌，创建系统访问记录
	token.Oauth2InfoCreate(&info, deviceFingerprint, [4]string{ipaddr, location, os, browser})
	s.oauth2Service.UpdateLoginDateAndIP(info)
	s.oauth2LogLoginService.Insert(
		info.ClientId, constants.STATUS_YES, "刷新访问令牌成功",
		[4]string{ipaddr, location, os, browser},
	)

	// 返回访问令牌和刷新令牌
	c.JSON(200, resp.OkData(map[string]any{
		"tokenType":        constants.HEADER_PREFIX,
		"accessToken":      accessToken,
		"expiresIn":        expiresIn,
		"refreshToken":     refreshToken,
		"refreshExpiresIn": refreshExpiresIn,
	}))
}
