package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/infra/cache"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/crypto"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/web/token"
	"mask_api_gin/src/modules/system/model/vo"
	systemService "mask_api_gin/src/modules/system/service"
)

// NewAccount 实例化服务层
var NewAccount = &Account{
	sysUserService:   systemService.NewSysUser,
	sysConfigService: systemService.NewSysConfig,
	sysRoleService:   systemService.NewSysRole,
	sysMenuService:   systemService.NewSysMenu,
}

// Account 账号身份操作 服务层处理
type Account struct {
	sysUserService   *systemService.SysUser   // 用户信息服务
	sysConfigService *systemService.SysConfig // 参数配置服务
	sysRoleService   *systemService.SysRole   // 角色服务
	sysMenuService   *systemService.SysMenu   // 菜单服务
}

// ValidateCaptcha 校验验证码
func (s Account) ValidateCaptcha(code, uuid string) error {
	// 验证码检查，从数据库配置获取验证码开关 true开启，false关闭
	captchaEnabledStr := s.sysConfigService.FindValueByKey("sys.account.captcha_enabled")
	if !parse.Boolean(captchaEnabledStr) {
		return nil
	}
	if code == "" || uuid == "" {
		return fmt.Errorf("captcha empty")
	}
	verifyKey := constants.CACHE_CAPTCHA_CODE + ":" + uuid
	captcha, _ := cache.GetAndDel("", verifyKey)
	if captcha == "" {
		return fmt.Errorf("captcha expire")
	}
	if captcha != strings.ToLower(code) {
		return fmt.Errorf("captcha error")
	}
	return nil
}

// ValidateCaptchaSMS 校验并一次性消费短信登录验证码。
func (s Account) ValidateCaptchaSMS(phone, code, uuid string) error {
	if len(code) != 6 {
		return fmt.Errorf("短信验证码格式错误")
	}
	if len(uuid) == 0 {
		return fmt.Errorf("短信验证码已失效")
	}

	codeKey := fmt.Sprintf("%s:%s", constants.CACHE_SMS_CODE, uuid)
	codeValue, _ := cache.Get("", codeKey)
	if codeValue == "" {
		return fmt.Errorf("短信验证码已失效")
	}
	expected := fmt.Sprintf("%s:%s", phone, code)
	if codeValue != expected {
		return fmt.Errorf("短信验证码错误")
	}
	// 验证码校验通过，删除缓存
	_ = cache.Del("", codeKey)
	// 验证码已被消费，同步清理发送冷却，允许用户立刻重新获取。
	_ = cache.Del("", fmt.Sprintf("%s:%s", constants.CACHE_SMS_CODE_COOLDOWN, phone))
	return nil
}

// ByUsername 账号密码登录创建用户信息
func (s Account) ByUsername(username, password string) (token.UserInfo, error) {
	info := token.UserInfo{}

	// 检查密码重试次数
	retryKey, _, lockTime, err := s.passwordRetryCount(username)
	if err != nil {
		return info, err
	}

	// 查询用户登录账号
	sysUser := s.sysUserService.FindByUserName(username)
	if sysUser.UserName != username {
		return info, fmt.Errorf("user does not exist or password is incorrect")
	}
	if sysUser.UserDel == constants.STATUS_YES {
		return info, fmt.Errorf("sorry, your account has been deleted. Sorry, your account has been deleted")
	}
	if sysUser.UserStatus == constants.STATUS_NO {
		return info, fmt.Errorf("sorry, your account has been disabled")
	}

	// 检验用户密码
	compareBool := crypto.Argon2idCompare(password, sysUser.UserPassword)
	if compareBool {
		s.CleanLoginRecordCache(sysUser.UserName) // 清除错误记录次数
	} else {
		_, _ = cache.IncrWithExpire("", retryKey, int64(lockTime.Seconds()))
		return info, fmt.Errorf("user does not exist or password is incorrect")
	}

	// 登录用户信息
	info.UserId = sysUser.ID
	info.DeptId = sysUser.DeptId
	info.UserName = sysUser.UserName
	roleKeys := make([]string, 0, len(sysUser.Roles))
	for _, r := range sysUser.Roles {
		roleKeys = append(roleKeys, r.RoleKey)
	}
	info.RoleKeys = roleKeys
	info.User, _ = json.Marshal(sysUser)
	// 用户权限组标识
	if config.IsSystemUser(sysUser.ID) {
		info.Permissions = []string{constants.SYS_PERMISSION_SYSTEM}
	} else {
		perms := s.sysMenuService.FindPermsByUserId(sysUser.ID)
		info.Permissions = parse.RemoveDuplicates(perms)
	}
	return info, nil
}

// ByPhone 手机号验证码登录创建用户信息。
func (s Account) ByPhone(phone string) (token.UserInfo, error) {
	info := token.UserInfo{}

	// 查询用户登录账号
	sysUser := s.sysUserService.FindByPhone(phone)
	if sysUser.UserPhone != phone || sysUser.ID <= 0 {
		return info, fmt.Errorf("user does not exist")
	}
	if sysUser.UserDel == constants.STATUS_YES {
		return info, fmt.Errorf("sorry, your account has been deleted. Sorry, your account has been deleted")
	}
	if sysUser.UserStatus == constants.STATUS_NO {
		return info, fmt.Errorf("sorry, your account has been disabled")
	}

	// 登录用户信息
	info.UserId = sysUser.ID
	info.DeptId = sysUser.DeptId
	info.UserName = sysUser.UserName
	roleKeys := make([]string, 0, len(sysUser.Roles))
	for _, r := range sysUser.Roles {
		roleKeys = append(roleKeys, r.RoleKey)
	}
	info.RoleKeys = roleKeys
	info.User, _ = json.Marshal(sysUser)
	// 用户权限组标识
	if config.IsSystemUser(sysUser.ID) {
		info.Permissions = []string{constants.SYS_PERMISSION_SYSTEM}
	} else {
		perms := s.sysMenuService.FindPermsByUserId(sysUser.ID)
		info.Permissions = parse.RemoveDuplicates(perms)
	}
	return info, nil
}

// ByUserId 用户ID刷新令牌创建用户信息
func (s Account) ByUserId(userId int64) (token.UserInfo, error) {
	info := token.UserInfo{}

	// 查询用户登录账号
	sysUser := s.sysUserService.FindById(userId)
	if sysUser.ID != userId {
		return info, fmt.Errorf("user does not exist")
	}
	if sysUser.UserDel == constants.STATUS_YES {
		return info, fmt.Errorf("sorry, your account has been deleted. Sorry, your account has been deleted")
	}
	if sysUser.UserStatus == constants.STATUS_NO {
		return info, fmt.Errorf("sorry, your account has been disabled")
	}

	// 登录用户信息
	info.UserId = sysUser.ID
	info.DeptId = sysUser.DeptId
	info.UserName = sysUser.UserName
	roleKeys := make([]string, 0, len(sysUser.Roles))
	for _, r := range sysUser.Roles {
		roleKeys = append(roleKeys, r.RoleKey)
	}
	info.RoleKeys = roleKeys
	info.User, _ = json.Marshal(sysUser)
	// 用户权限组标识
	if config.IsSystemUser(sysUser.ID) {
		info.Permissions = []string{constants.SYS_PERMISSION_SYSTEM}
	} else {
		perms := s.sysMenuService.FindPermsByUserId(sysUser.ID)
		info.Permissions = parse.RemoveDuplicates(perms)
	}
	return info, nil
}

// UpdateLoginDateAndIP 更新登录时间和IP
func (s Account) UpdateLoginDateAndIP(info token.UserInfo) bool {
	user := s.sysUserService.FindById(info.UserId)
	user.UserPassword = "" // 密码不更新
	user.UserLoginIp = info.LoginIp
	user.UserLoginTime = info.LoginTime
	return s.sysUserService.Update(user) > 0
}

// CleanLoginRecordCache 清除错误记录次数
func (s Account) CleanLoginRecordCache(userName string) bool {
	base64UserName := base64.StdEncoding.EncodeToString([]byte(userName))
	cacheKey := fmt.Sprintf("%s:%s", constants.CACHE_PWD_ERR_COUNT, base64UserName)
	return cache.Del("", cacheKey) == nil
}

// passwordRetryCount 密码重试次数
func (s Account) passwordRetryCount(userName string) (string, int64, time.Duration, error) {
	// 验证登录次数和错误锁定时间
	passwordRetryMaxCount := s.sysConfigService.FindValueByKey("sys.user.max_retry_count")
	retryMaxCount := parse.Number(passwordRetryMaxCount)
	if retryMaxCount <= 0 {
		retryMaxCount = 3
	}
	passwordRetryLockTime := s.sysConfigService.FindValueByKey("sys.user.lock_time")
	retryLockTime := parse.Number(passwordRetryLockTime)
	if retryLockTime <= 0 {
		retryLockTime = 1
	}
	// 验证缓存记录次数
	base64UserName := base64.StdEncoding.EncodeToString([]byte(userName))
	retryKey := fmt.Sprintf("%s:%s", constants.CACHE_PWD_ERR_COUNT, base64UserName)
	retryCount, err := cache.Get("", retryKey)
	if retryCount == "" || err != nil {
		retryCount = "0"
	}
	// 是否超过错误值
	retryCountInt64 := parse.Number(retryCount)
	if retryCountInt64 >= retryMaxCount {
		msg := fmt.Sprintf("密码输入错误 %d 次，帐户锁定 %d 分钟", retryMaxCount, retryLockTime)
		return retryKey, retryCountInt64, time.Duration(retryLockTime) * time.Minute, fmt.Errorf("%s", msg)
	}
	return retryKey, retryCountInt64, time.Duration(retryLockTime) * time.Minute, nil
}

// RoleAndMenuPerms 角色和菜单数据权限
func (s Account) RoleAndMenuPerms(userId int64, isSystemUser bool) ([]string, []string) {
	if isSystemUser {
		return []string{constants.SYS_ROLE_SYSTEM_KEY}, []string{constants.SYS_PERMISSION_SYSTEM}
	}
	// 角色key
	var roleGroup []string
	roles := s.sysRoleService.FindByUserId(userId)
	for _, role := range roles {
		roleGroup = append(roleGroup, role.RoleKey)
	}
	// 菜单权限key
	perms := s.sysMenuService.FindPermsByUserId(userId)
	return parse.RemoveDuplicates(roleGroup), parse.RemoveDuplicates(perms)
}

// RouteMenus 前端路由所需要的菜单
func (s Account) RouteMenus(userId int64, isSystemUser bool) []vo.Router {
	var buildMenus []vo.Router
	if isSystemUser {
		menus := s.sysMenuService.BuildTreeMenusByUserId(0)
		buildMenus = s.sysMenuService.BuildRouteMenus(menus, "")
	} else {
		menus := s.sysMenuService.BuildTreeMenusByUserId(userId)
		buildMenus = s.sysMenuService.BuildRouteMenus(menus, "")
	}
	return buildMenus
}
