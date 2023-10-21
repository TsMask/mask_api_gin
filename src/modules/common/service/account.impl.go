package service

import (
	"errors"
	"fmt"
	"mask_api_gin/src/framework/config"
	adminConstants "mask_api_gin/src/framework/constants/admin"
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/utils/crypto"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
	systemService "mask_api_gin/src/modules/system/service"
	"time"
)

// 实例化服务层 AccountImpl 结构体
var NewAccountImpl = &AccountImpl{
	sysUserService:   systemService.NewSysUserImpl,
	sysConfigService: systemService.NewSysConfigImpl,
	sysRoleService:   systemService.NewSysRoleImpl,
	sysMenuService:   systemService.NewSysMenuImpl,
}

// 账号身份操作服务 服务层处理
type AccountImpl struct {
	// 用户信息服务
	sysUserService systemService.ISysUser
	// 参数配置服务
	sysConfigService systemService.ISysConfig
	// 角色服务
	sysRoleService systemService.ISysRole
	// 菜单服务
	sysMenuService systemService.ISysMenu
}

// ValidateCaptcha 校验验证码
func (s *AccountImpl) ValidateCaptcha(code, uuid string) error {
	// 验证码检查，从数据库配置获取验证码开关 true开启，false关闭
	captchaEnabledStr := s.sysConfigService.SelectConfigValueByKey("sys.account.captchaEnabled")
	if !parse.Boolean(captchaEnabledStr) {
		return nil
	}
	if code == "" || uuid == "" {
		return errors.New("验证码信息错误")
	}
	verifyKey := cachekey.CAPTCHA_CODE_KEY + uuid
	captcha, _ := redis.Get("", verifyKey)
	if captcha == "" {
		return errors.New("验证码已失效")
	}
	redis.Del("", verifyKey)
	if captcha != code {
		return errors.New("验证码错误")
	}
	return nil
}

// LoginByUsername 登录创建用户信息
func (s *AccountImpl) LoginByUsername(username, password string) (vo.LoginUser, error) {
	loginUser := vo.LoginUser{}

	// 检查密码重试次数
	retrykey, retryCount, lockTime, err := s.passwordRetryCount(username)
	if err != nil {
		return loginUser, err
	}

	// 查询用户登录账号
	sysUser := s.sysUserService.SelectUserByUserName(username)
	if sysUser.UserName != username {
		return loginUser, errors.New("用户不存在或密码错误")
	}
	if sysUser.DelFlag == common.STATUS_YES {
		return loginUser, errors.New("对不起，您的账号已被删除")
	}
	if sysUser.Status == common.STATUS_NO {
		return loginUser, errors.New("对不起，您的账号已禁用")
	}

	// 检验用户密码
	compareBool := crypto.BcryptCompare(password, sysUser.Password)
	if !compareBool {
		redis.SetByExpire("", retrykey, retryCount+1, lockTime)
		return loginUser, errors.New("用户不存在或密码错误")
	} else {
		// 清除错误记录次数
		s.ClearLoginRecordCache(username)
	}

	// 登录用户信息
	loginUser.UserID = sysUser.UserID
	loginUser.DeptID = sysUser.DeptID
	loginUser.User = sysUser
	// 用户权限组标识
	isAdmin := config.IsAdmin(sysUser.UserID)
	if isAdmin {
		loginUser.Permissions = []string{adminConstants.PERMISSION}
	} else {
		perms := s.sysMenuService.SelectMenuPermsByUserId(sysUser.UserID)
		loginUser.Permissions = parse.RemoveDuplicates(perms)
	}
	return loginUser, nil
}

// UpdateLoginDateAndIP 更新登录时间和IP
func (s *AccountImpl) UpdateLoginDateAndIP(loginUser *vo.LoginUser) bool {
	sysUser := loginUser.User
	userInfo := model.SysUser{
		UserID:    sysUser.UserID,
		LoginIP:   sysUser.LoginIP,
		LoginDate: sysUser.LoginDate,
		UpdateBy:  sysUser.UserName,
	}
	rows := s.sysUserService.UpdateUser(userInfo)
	return rows > 0
}

// ClearLoginRecordCache 清除错误记录次数
func (s *AccountImpl) ClearLoginRecordCache(username string) bool {
	cacheKey := cachekey.PWD_ERR_CNT_KEY + username
	hasKey, _ := redis.Has("", cacheKey)
	if hasKey {
		delOk, _ := redis.Del("", cacheKey)
		return delOk
	}
	return false
}

// passwordRetryCount 密码重试次数
func (s *AccountImpl) passwordRetryCount(username string) (string, int64, time.Duration, error) {
	// 验证登录次数和错误锁定时间
	maxRetryCount := config.Get("user.password.maxRetryCount").(int)
	lockTime := config.Get("user.password.lockTime").(int)
	// 验证缓存记录次数
	retrykey := cachekey.PWD_ERR_CNT_KEY + username
	retryCount, err := redis.Get("", retrykey)
	if retryCount == "" || err != nil {
		retryCount = "0"
	}
	// 是否超过错误值
	retryCountInt64 := parse.Number(retryCount)
	if retryCountInt64 >= int64(maxRetryCount) {
		msg := fmt.Sprintf("密码输入错误 %d 次，帐户锁定 %d 分钟", maxRetryCount, lockTime)
		return retrykey, retryCountInt64, time.Duration(lockTime) * time.Minute, errors.New(msg)
	}
	return retrykey, retryCountInt64, time.Duration(lockTime) * time.Minute, nil
}

// RoleAndMenuPerms 角色和菜单数据权限
func (s *AccountImpl) RoleAndMenuPerms(userId string, isAdmin bool) ([]string, []string) {
	if isAdmin {
		return []string{adminConstants.ROLE_KEY}, []string{adminConstants.PERMISSION}
	} else {
		// 角色key
		roleGroup := []string{}
		roles := s.sysRoleService.SelectRoleListByUserId(userId)
		for _, role := range roles {
			roleGroup = append(roleGroup, role.RoleKey)
		}
		// 菜单权限key
		perms := s.sysMenuService.SelectMenuPermsByUserId(userId)
		return parse.RemoveDuplicates(roleGroup), parse.RemoveDuplicates(perms)
	}
}

// RouteMenus 前端路由所需要的菜单
func (s *AccountImpl) RouteMenus(userId string, isAdmin bool) []vo.Router {
	var buildMenus []vo.Router
	if isAdmin {
		menus := s.sysMenuService.SelectMenuTreeByUserId("*")
		buildMenus = s.sysMenuService.BuildRouteMenus(menus, "")
	} else {
		menus := s.sysMenuService.SelectMenuTreeByUserId(userId)
		buildMenus = s.sysMenuService.BuildRouteMenus(menus, "")
	}
	return buildMenus
}
