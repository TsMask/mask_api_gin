package service

import (
	"fmt"
	"strings"
	"time"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/infra/cache"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/generate"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/shared/regular"
	"mask_api_gin/src/internal/web/token"
	"mask_api_gin/src/modules/auth/model"
	systemService "mask_api_gin/src/modules/system/service"
)

// NewPasswordReset 实例化服务层
var NewPasswordReset = &PasswordReset{
	sysUserService:   systemService.NewSysUser,
	sysConfigService: systemService.NewSysConfig,
}

// PasswordReset 找回密码操作 服务层处理
type PasswordReset struct {
	sysUserService   *systemService.SysUser   // 用户信息服务
	sysConfigService *systemService.SysConfig // 参数配置服务
}

// Lookup 查询密码找回联系方式
func (s PasswordReset) Lookup(username string) (map[string]string, error) {
	target, err := s.findTarget(username)
	if err != nil {
		return nil, err
	}
	maskedContact := maskPhone(target.Contact)
	if target.Channel == "email" {
		maskedContact = maskEmail(target.Contact)
	}
	return map[string]string{
		"channel":        target.Channel,
		"masked_contact": maskedContact,
	}, nil
}

// SendCaptchaEmail 发送找回密码验证码。本地环境返回验证码便于调试。
func (s PasswordReset) SendCaptchaEmail(username string) (map[string]any, error) {
	target, err := s.findTarget(username)
	if err != nil {
		return nil, err
	}
	if config.Env() != "local" {
		if target.Channel == "email" {
			return nil, fmt.Errorf("邮箱服务未配置")
		}
		return nil, fmt.Errorf("短信服务未配置")
	}

	expiresIn := int64(300)
	resendInterval := int64(60)
	if target.Channel == "email" {
		expiresIn = parse.Number(s.sysConfigService.FindValueByKey("sys.account.captchaEmailExpiresIn"))
		resendInterval = parse.Number(s.sysConfigService.FindValueByKey("sys.account.captchaEmailResendInterval"))
	} else {
		expiresIn = parse.Number(s.sysConfigService.FindValueByKey("sys.account.captchaSMSExpiresIn"))
		resendInterval = parse.Number(s.sysConfigService.FindValueByKey("sys.account.captchaSMSResendInterval"))
	}
	if expiresIn <= 0 {
		expiresIn = 300
	}
	if resendInterval <= 0 {
		resendInterval = 60
	}

	cooldownKey := fmt.Sprintf("%s:%s", constants.CACHE_PASSWORD_RESET_COOLDOWN, target.Contact)
	if value, _ := cache.Get("", cooldownKey); value != "" {
		return nil, fmt.Errorf("验证码发送过于频繁，请稍后再试")
	}

	code := fmt.Sprintf("%06d", generate.Number(6))
	uuid := generate.String(32)
	codeKey := fmt.Sprintf("%s:%s", constants.CACHE_PASSWORD_RESET_CODE, uuid)
	codeValue := fmt.Sprintf("%s:%s:%s", target.UserName, target.Contact, code)
	if err := cache.Set("", codeKey, codeValue, time.Duration(expiresIn)*time.Second); err != nil {
		return nil, fmt.Errorf("验证码生成失败")
	}
	if err := cache.Set("", cooldownKey, target.Contact, time.Duration(resendInterval)*time.Second); err != nil {
		_ = cache.Del("", codeKey)
		return nil, fmt.Errorf("验证码生成失败")
	}

	data := map[string]any{
		"uuid":            uuid,
		"expires_in":      expiresIn,
		"resend_interval": resendInterval,
	}
	if config.Env() == "local" {
		data["text"] = code
	}
	return data, nil
}

// Reset 校验找回验证码并重置密码。
func (s PasswordReset) Reset(username, code, uuid, password string) error {
	target, err := s.findTarget(username)
	if err != nil {
		return err
	}
	if len(code) != 6 || uuid == "" {
		return fmt.Errorf("验证码格式错误")
	}

	codeKey := fmt.Sprintf("%s:%s", constants.CACHE_PASSWORD_RESET_CODE, uuid)
	codeValue, _ := cache.Get("", codeKey)
	if codeValue == "" {
		return fmt.Errorf("验证码已失效")
	}
	expected := fmt.Sprintf("%s:%s:%s", target.UserName, target.Contact, code)
	if codeValue != expected {
		return fmt.Errorf("验证码错误")
	}

	user := s.sysUserService.FindByUserName(target.UserName)
	if user.ID != target.UserID {
		return fmt.Errorf("用户不存在或账号不可自助找回密码")
	}
	user.UserPassword = password
	user.UpdateBy = username
	if s.sysUserService.Update(user) <= 0 {
		return fmt.Errorf("重置密码失败，请联系系统管理员")
	}

	// 密码重置后清理验证码、冷却和登录态，强制使用新密码登录。
	_ = cache.Del("", codeKey)
	_ = cache.Del("", fmt.Sprintf("%s:%s", constants.CACHE_PASSWORD_RESET_COOLDOWN, target.Contact))
	token.RemoveByUserIds([]int64{target.UserID})
	return nil
}

// findTarget 查询用户并确定找回联系方式。
func (s PasswordReset) findTarget(username string) (model.PasswordResetTarget, error) {
	if !regular.ValidUsername(username) {
		return model.PasswordResetTarget{}, fmt.Errorf("用户账号格式不正确")
	}
	user := s.sysUserService.FindByUserName(username)
	if user.UserName != username || user.ID <= 0 || user.UserDel == constants.STATUS_YES || user.UserStatus == constants.STATUS_NO {
		return model.PasswordResetTarget{}, fmt.Errorf("用户不存在或账号不可自助找回密码")
	}
	if regular.ValidEmail(user.UserEmail) {
		return model.PasswordResetTarget{
			UserID:   user.ID,
			UserName: user.UserName,
			Channel:  "email",
			Contact:  user.UserEmail,
		}, nil
	}
	if regular.ValidMobile(user.UserPhone) {
		return model.PasswordResetTarget{
			UserID:   user.ID,
			UserName: user.UserName,
			Channel:  "phone",
			Contact:  user.UserPhone,
		}, nil
	}
	return model.PasswordResetTarget{}, fmt.Errorf("账号未绑定邮箱或手机号，请联系系统管理员")
}

// maskEmail 邮箱地址脱敏
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" {
		return "已绑定邮箱"
	}
	local := []rune(parts[0])
	return fmt.Sprintf("%c***@%s", local[0], parts[1])
}

// maskPhone 手机号脱敏
func maskPhone(phone string) string {
	if len(phone) != 11 {
		return "已绑定手机号"
	}
	return fmt.Sprintf("%s****%s", phone[:3], phone[7:])
}
