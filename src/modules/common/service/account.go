package service

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	constCacheKey "mask_api_gin/src/framework/constants/cache_key"
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/utils/crypto"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo"
	systemService "mask_api_gin/src/modules/system/service"
	"time"
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
	captchaEnabledStr := s.sysConfigService.FindValueByKey("sys.account.captchaEnabled")
	if !parse.Boolean(captchaEnabledStr) {
		return nil
	}
	if code == "" || uuid == "" {
		return fmt.Errorf("验证码信息错误")
	}
	verifyKey := constCacheKey.CAPTCHA_CODE_KEY + uuid
	captcha, _ := redis.Get("", verifyKey)
	if captcha == "" {
		return fmt.Errorf("验证码已失效")
	}
	_ = redis.Del("", verifyKey)
	if captcha != code {
		return fmt.Errorf("验证码错误")
	}
	return nil
}

// ByUsername 登录创建用户信息
func (s Account) ByUsername(username, password string) (vo.LoginUser, error) {
	var loginUser vo.LoginUser

	// 检查密码重试次数
	retryKey, retryCount, lockTime, err := s.passwordRetryCount(username)
	if err != nil {
		return loginUser, err
	}

	// 查询用户登录账号
	sysUser := s.sysUserService.FindByUserName(username)
	if sysUser.UserName != username {
		return loginUser, fmt.Errorf("用户不存在或密码错误")
	}
	if sysUser.DelFlag == constSystem.STATUS_YES {
		return loginUser, fmt.Errorf("对不起，您的账号已被删除")
	}
	if sysUser.Status == constSystem.STATUS_NO {
		return loginUser, fmt.Errorf("对不起，您的账号已禁用")
	}

	// 检验用户密码
	compareBool := crypto.BcryptCompare(password, sysUser.Password)
	if compareBool {
		s.CleanLoginRecordCache(username) // 清除错误记录次数
	} else {
		_ = redis.SetByExpire("", retryKey, retryCount+1, lockTime)
		return loginUser, fmt.Errorf("用户不存在或密码错误")
	}

	// 登录用户信息
	loginUser = vo.LoginUser{}
	loginUser.UserID = sysUser.UserId
	loginUser.DeptID = sysUser.DeptId
	loginUser.User = sysUser
	// 用户权限组标识
	if config.IsSysAdmin(sysUser.UserId) {
		loginUser.Permissions = []string{constSystem.PERMISSION_SYSTEM}
	} else {
		perms := s.sysMenuService.FindPermsByUserId(sysUser.UserId)
		loginUser.Permissions = parse.RemoveDuplicates(perms)
	}
	return loginUser, nil
}

// UpdateLoginDateAndIP 更新登录时间和IP
func (s Account) UpdateLoginDateAndIP(loginUser *vo.LoginUser) bool {
	sysUser := loginUser.User
	user := s.sysUserService.FindById(sysUser.UserId)
	user.Password = "" // 密码不更新
	user.LoginIp = sysUser.LoginIp
	user.LoginDate = sysUser.LoginDate
	return s.sysUserService.Update(user) > 0
}

// CleanLoginRecordCache 清除错误记录次数
func (s Account) CleanLoginRecordCache(username string) bool {
	cacheKey := constCacheKey.PWD_ERR_COUNT_KEY + username
	hasKey, err := redis.Has("", cacheKey)
	if hasKey > 0 && err == nil {
		return redis.Del("", cacheKey) == nil
	}
	return false
}

// passwordRetryCount 密码重试次数
func (s Account) passwordRetryCount(username string) (string, int64, time.Duration, error) {
	// 验证登录次数和错误锁定时间
	maxRetryCount := config.Get("user.password.maxRetryCount").(int)
	lockTime := config.Get("user.password.lockTime").(int)
	// 验证缓存记录次数
	retryKey := constCacheKey.PWD_ERR_COUNT_KEY + username
	retryCount, err := redis.Get("", retryKey)
	if retryCount == "" || err != nil {
		retryCount = "0"
	}
	// 是否超过错误值
	retryCountInt64 := parse.Number(retryCount)
	if retryCountInt64 >= int64(maxRetryCount) {
		msg := fmt.Sprintf("密码输入错误 %d 次，帐户锁定 %d 分钟", maxRetryCount, lockTime)
		return retryKey, retryCountInt64, time.Duration(lockTime) * time.Minute, fmt.Errorf("%s", msg)
	}
	return retryKey, retryCountInt64, time.Duration(lockTime) * time.Minute, nil
}

// RoleAndMenuPerms 角色和菜单数据权限
func (s Account) RoleAndMenuPerms(userId string, isSysAdmin bool) ([]string, []string) {
	if isSysAdmin {
		return []string{constSystem.ROLE_SYSTEM_KEY}, []string{constSystem.PERMISSION_SYSTEM}
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
func (s Account) RouteMenus(userId string, isSysAdmin bool) []vo.Router {
	var buildMenus []vo.Router
	if isSysAdmin {
		menus := s.sysMenuService.BuildTreeMenusByUserId("*")
		buildMenus = s.sysMenuService.BuildRouteMenus(menus, "")
	} else {
		menus := s.sysMenuService.BuildTreeMenusByUserId(userId)
		buildMenus = s.sysMenuService.BuildRouteMenus(menus, "")
	}
	return buildMenus
}
