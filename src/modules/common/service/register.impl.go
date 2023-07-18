package service

import (
	"errors"
	"fmt"
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/utils/parse"
	monitorService "mask_api_gin/src/modules/monitor/service"
	systemModel "mask_api_gin/src/modules/system/model"
	systemService "mask_api_gin/src/modules/system/service"
)

// 账号身份操作服务 业务层处理
var RegisterImpl = &registerImpl{
	sysLogininforService: monitorService.SysLogininforImpl,
	sysUserService:       systemService.SysUserImpl,
	sysConfigService:     systemService.SysConfigImpl,
	sysRoleService:       systemService.SysRoleImpl,
}

type registerImpl struct {
	// 系统登录访问信息服务
	sysLogininforService monitorService.ISysLogininfor
	// 用户信息服务
	sysUserService systemService.ISysUser
	// 参数配置服务
	sysConfigService systemService.ISysConfig
	// 角色服务
	sysRoleService systemService.ISysRole
}

// ValidateCaptcha 校验验证码
func (s *registerImpl) ValidateCaptcha(code, uuid string) error {
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

// ByUserName 账号注册
func (s *registerImpl) ByUserName(username, password, userType string) string {
	// 检查用户登录账号是否唯一
	uniqueUserName := s.sysUserService.CheckUniqueUserName(username, "")
	if !uniqueUserName {
		return fmt.Sprintf("注册用户【%s】失败，注册账号已存在", username)
	}

	sysUser := systemModel.SysUser{
		UserName: username,
		NickName: username,          // 昵称使用名称账号
		Password: password,          // 原始密码
		Status:   common.STATUS_YES, // 账号状态激活
		DeptID:   "100",             // 归属部门为根节点
		CreateBy: "注册",              // 创建来源
	}
	// 标记用户类型
	if userType == "" {
		sysUser.UserType = "sys"
	}
	// 新增用户的角色管理
	sysUser.RoleIDs = s.registerRoleInit(userType)
	// 新增用户的岗位管理
	sysUser.PostIDs = s.registerPostInit(userType)

	insertId := s.sysUserService.InsertUser(sysUser)
	if insertId != "" {
		return insertId
	}
	return "注册失败，请联系系统管理人员"
}

// registerRoleInit 注册初始角色
func (s *registerImpl) registerRoleInit(userType string) []string {
	if userType == "sys" {
		return []string{}
	}
	return []string{}
}

// registerPostInit 注册初始岗位
func (s *registerImpl) registerPostInit(userType string) []string {
	if userType == "sys" {
		return []string{}
	}
	return []string{}
}