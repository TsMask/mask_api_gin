package service

import (
	"errors"
	"fmt"
	monitorService "mask_api_gin/src/modules/monitor/service"
	systemService "mask_api_gin/src/modules/system/service"
	"mask_api_gin/src/pkg/cache/redis"
	"mask_api_gin/src/pkg/config"
	adminConstants "mask_api_gin/src/pkg/constants/admin"
	"mask_api_gin/src/pkg/constants/cachekey"
	"mask_api_gin/src/pkg/constants/common"
	"mask_api_gin/src/pkg/model"
	"mask_api_gin/src/pkg/utils/crypto"
	"mask_api_gin/src/pkg/utils/parse"
	"time"
)

// 账号身份操作服务 业务层处理
var AccountImpl = &accountImpl{
	sysLogininforService: monitorService.SysLogininforImpl,
	sysUserService:       systemService.SysUserImpl,
	sysConfigService:     systemService.SysConfigImpl,
	sysRoleService:       systemService.SysRoleImpl,
	sysMenuService:       systemService.SysMenuImpl,
}

type accountImpl struct {
	// 系统登录访问信息服务
	sysLogininforService monitorService.ISysLogininfor
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
func (s *accountImpl) ValidateCaptcha(username, code, uuid string) error {
	// 验证码检查，从数据库配置获取验证码开关 true开启，false关闭
	captchaEnabledStr := s.sysConfigService.SelectConfigValueByKey("sys.account.captchaEnabled")
	if parse.Boolean(captchaEnabledStr) {
		verifyKey := cachekey.CAPTCHA_CODE_KEY + uuid
		captcha := redis.Get(verifyKey)
		if captcha == "" {
			return errors.New("验证码已失效")
		}
		redis.Del(verifyKey)
		if captcha != code {
			return errors.New("验证码错误")
		}
	}
	return nil
}

// LoginByUsername 登录创建用户信息
func (s *accountImpl) LoginByUsername(username, password string) (model.LoginUser, error) {
	// 检查密码重试次数
	retrykey, retryCount, lockTime, err := s.passwordRetryCount(username)
	if err != nil {
		return model.LoginUser{}, err
	}

	// 查询用户登录账号
	sysUser := s.sysUserService.SelectUserByUserName(username)
	if sysUser.UserName != username {
		return model.LoginUser{}, errors.New("用户不存在或密码错误")
	}
	if sysUser.DelFlag == common.STATUS_YES {
		return model.LoginUser{}, errors.New("对不起，您的账号已被删除")
	}
	if sysUser.Status == common.STATUS_NO {
		return model.LoginUser{}, errors.New("用户已封禁，请联系管理员")
	}

	// 检验用户密码
	compareBool := crypto.BcryptCompare(password, sysUser.Password)
	if !compareBool {
		redis.SetByExpire(retrykey, retryCount+1, lockTime)
		return model.LoginUser{}, errors.New("用户不存在或密码错误")
	} else {
		// 清除错误记录次数
		s.ClearLoginRecordCache(username)
	}

	// 登录用户信息
	loginUser := model.LoginUser{
		UserID: sysUser.UserID,
		DeptID: sysUser.DeptID,
		User:   sysUser,
	}
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

// ClearLoginRecordCache 清除错误记录次数
func (s *accountImpl) ClearLoginRecordCache(loginName string) bool {
	cacheKey := cachekey.PWD_ERR_CNT_KEY + loginName
	if redis.Has(cacheKey) {
		return redis.Del(cacheKey)
	}
	return false
}

// passwordRetryCount 密码重试次数
func (s *accountImpl) passwordRetryCount(username string) (string, int64, time.Duration, error) {
	// 验证登录次数和错误锁定时间
	maxRetryCount := config.Get("user.password.maxRetryCount").(int)
	lockTime := config.Get("user.password.lockTime").(int)
	// 验证缓存记录次数
	retrykey := cachekey.PWD_ERR_CNT_KEY + username
	retryCount := redis.Get(retrykey)
	if retryCount == "" {
		retryCount = "0"
	}
	// 是否超过错误值
	if parse.Number(retryCount) >= int64(maxRetryCount) {
		msg := fmt.Sprintf("密码输入错误 %d 次，帐户锁定 %d 分钟", maxRetryCount, lockTime)
		return retrykey, int64(maxRetryCount), time.Duration(lockTime) * time.Minute, errors.New(msg)
	}
	return retrykey, int64(maxRetryCount), time.Duration(lockTime) * time.Minute, nil
}

// RoleAndMenuPerms 角色和菜单数据权限 TODO
func (s *accountImpl) RoleAndMenuPerms(userId string, isAdmin bool) ([]string, []string) {
	if isAdmin {
		return []string{adminConstants.ROLE_KEY}, []string{adminConstants.PERMISSION}
	} else {
		roles := s.sysRoleService.SelectRolePermsByUserId(userId)
		perms := s.sysMenuService.SelectMenuPermsByUserId(userId)
		return parse.RemoveDuplicates(roles), parse.RemoveDuplicates(perms)
	}
}

// RouteMenus 前端路由所需要的菜单 TODO
func (s *accountImpl) RouteMenus(userId string, isAdmin bool) []model.Router {
	var buildMenus []model.Router
	if isAdmin {
		menus := s.sysMenuService.SelectMenuTreeByUserId("")
		buildMenus = s.sysMenuService.BuildRouteMenus(menus)
	} else {
		menus := s.sysMenuService.SelectMenuTreeByUserId(userId)
		buildMenus = s.sysMenuService.BuildRouteMenus(menus)
	}
	return buildMenus
}
