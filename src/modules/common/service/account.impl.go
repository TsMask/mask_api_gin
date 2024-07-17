package service

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	constAdmin "mask_api_gin/src/framework/constants/admin"
	constCacheKey "mask_api_gin/src/framework/constants/cache_key"
	constCommon "mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/utils/crypto"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo"
	systemService "mask_api_gin/src/modules/system/service"
	"time"
)

// NewAccountService 实例化服务层
var NewAccountService = &AccountServiceImpl{
	sysUserService:   systemService.NewSysUser,
	sysConfigService: systemService.NewSysConfig,
	sysRoleService:   systemService.NewSysRole,
	sysMenuService:   systemService.NewSysMenu,
}

// AccountServiceImpl 账号身份操作 服务层处理
type AccountServiceImpl struct {
	sysUserService   systemService.ISysUserService   // 用户信息服务
	sysConfigService systemService.ISysConfigService // 参数配置服务
	sysRoleService   systemService.ISysRoleService   // 角色服务
	sysMenuService   systemService.ISysMenuService   // 菜单服务
}

// ValidateCaptcha 校验验证码
func (s *AccountServiceImpl) ValidateCaptcha(code, uuid string) error {
	// 验证码检查，从数据库配置获取验证码开关 true开启，false关闭
	captchaEnabledStr := s.sysConfigService.FindValueByKey("sys.account.captchaEnabled")
	if !parse.Boolean(captchaEnabledStr) {
		return nil
	}
	if code == "" || uuid == "" {
		return fmt.Errorf("验证码信息错误")
	}
	verifyKey := constCacheKey.CaptchaCodeKey + uuid
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
func (s *AccountServiceImpl) ByUsername(username, password string) (vo.LoginUser, error) {
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
	if sysUser.DelFlag == constCommon.StatusYes {
		return loginUser, fmt.Errorf("对不起，您的账号已被删除")
	}
	if sysUser.Status == constCommon.StatusNo {
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
	loginUser.UserID = sysUser.UserID
	loginUser.DeptID = sysUser.DeptID
	loginUser.User = sysUser
	// 用户权限组标识
	isAdmin := config.IsAdmin(sysUser.UserID)
	if isAdmin {
		loginUser.Permissions = []string{constAdmin.Permission}
	} else {
		perms := s.sysMenuService.FindPermsByUserId(sysUser.UserID)
		loginUser.Permissions = parse.RemoveDuplicates(perms)
	}
	return loginUser, nil
}

// UpdateLoginDateAndIP 更新登录时间和IP
func (s *AccountServiceImpl) UpdateLoginDateAndIP(loginUser *vo.LoginUser) bool {
	sysUser := loginUser.User
	user := s.sysUserService.FindById(sysUser.UserID)
	user.LoginIP = sysUser.LoginIP
	user.LoginDate = sysUser.LoginDate
	return s.sysUserService.Update(user) > 0
}

// CleanLoginRecordCache 清除错误记录次数
func (s *AccountServiceImpl) CleanLoginRecordCache(username string) bool {
	cacheKey := constCacheKey.PwdErrCntKey + username
	hasKey, err := redis.Has("", cacheKey)
	if hasKey > 0 && err == nil {
		return redis.Del("", cacheKey) == nil
	}
	return false
}

// passwordRetryCount 密码重试次数
func (s *AccountServiceImpl) passwordRetryCount(username string) (string, int64, time.Duration, error) {
	// 验证登录次数和错误锁定时间
	maxRetryCount := config.Get("user.password.maxRetryCount").(int)
	lockTime := config.Get("user.password.lockTime").(int)
	// 验证缓存记录次数
	retryKey := constCacheKey.PwdErrCntKey + username
	retryCount, err := redis.Get("", retryKey)
	if retryCount == "" || err != nil {
		retryCount = "0"
	}
	// 是否超过错误值
	retryCountInt64 := parse.Number(retryCount)
	if retryCountInt64 >= int64(maxRetryCount) {
		msg := fmt.Sprintf("密码输入错误 %d 次，帐户锁定 %d 分钟", maxRetryCount, lockTime)
		return retryKey, retryCountInt64, time.Duration(lockTime) * time.Minute, fmt.Errorf(msg)
	}
	return retryKey, retryCountInt64, time.Duration(lockTime) * time.Minute, nil
}

// RoleAndMenuPerms 角色和菜单数据权限
func (s *AccountServiceImpl) RoleAndMenuPerms(userId string, isAdmin bool) ([]string, []string) {
	if isAdmin {
		return []string{constAdmin.RoleKey}, []string{constAdmin.Permission}
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
func (s *AccountServiceImpl) RouteMenus(userId string, isAdmin bool) []vo.Router {
	var buildMenus []vo.Router
	if isAdmin {
		menus := s.sysMenuService.BuildTreeMenusByUserId("*")
		buildMenus = s.sysMenuService.BuildRouteMenus(menus, "")
	} else {
		menus := s.sysMenuService.BuildTreeMenusByUserId(userId)
		buildMenus = s.sysMenuService.BuildRouteMenus(menus, "")
	}
	return buildMenus
}
