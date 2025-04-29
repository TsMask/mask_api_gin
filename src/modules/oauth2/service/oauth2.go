package service

import (
	"fmt"
	"strings"
	"time"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/token"
	"mask_api_gin/src/framework/utils/crypto"
	"mask_api_gin/src/framework/utils/generate"
	"mask_api_gin/src/modules/oauth2/model"
	"mask_api_gin/src/modules/oauth2/repository"
)

// NewOauth2Service 实例化服务层
var NewOauth2Service = &Oauth2Service{
	oauth2ClientRepository: repository.NewOauth2Client,
}

// Oauth2Service 用户授权第三方应用信息 服务层处理
type Oauth2Service struct {
	oauth2ClientRepository *repository.Oauth2Client // 用户授权第三方应用表
}

// CreateCode 创建授权码
func (s Oauth2Service) CreateCode() string {
	code := generate.Code(8)
	uuid := crypto.MD5(code)
	verifyKey := constants.CACHE_OAUTH2_CODE + ":" + uuid
	// 授权码有效期，单位秒
	codeExpiration := 2 * 60 * time.Second
	_ = redis.Set("", verifyKey, code, codeExpiration)
	return code
}

// ValidateCode 校验授权码
func (s Oauth2Service) ValidateCode(code string) error {
	if len(code) > 16 {
		return fmt.Errorf("code length error")
	}
	uuid := crypto.MD5(code)
	verifyKey := constants.CACHE_OAUTH2_CODE + ":" + uuid
	captcha, _ := redis.Get("", verifyKey)
	if captcha == "" {
		return fmt.Errorf("code expire")
	}
	_ = redis.Del("", verifyKey)
	if captcha != strings.ToLower(code) {
		return fmt.Errorf("code error")
	}
	return nil
}

// ByClient 客户端信息
func (s Oauth2Service) ByClient(clientId, clientSecret, ipaddr string) (token.Oauth2Info, error) {
	info := token.Oauth2Info{}

	// 查询用户登录账号
	var item model.Oauth2Client
	rows := s.oauth2ClientRepository.Select(model.Oauth2Client{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	})
	if len(rows) > 0 {
		item = rows[0]
	}
	if item.ClientId == "" || item.ClientSecret == "" {
		return info, fmt.Errorf("clientId or clientSecret is not exist")
	}
	// 判断IP白名单
	if !strings.Contains(item.IPWhite, ipaddr) {
		return info, fmt.Errorf("ip whitelist mismatch")
	}

	info.ClientId = clientId
	// 用户权限组标识
	info.Scope = []string{}
	return info, nil
}

// UpdateLoginDateAndIP 更新登录时间和IP
func (s Oauth2Service) UpdateLoginDateAndIP(info token.Oauth2Info) bool {
	item := s.oauth2ClientRepository.SelectByClientId(info.ClientId)
	item.LoginIp = info.LoginIp
	item.LoginTime = info.LoginTime
	rows := s.oauth2ClientRepository.Update(item)
	return rows > 0
}
