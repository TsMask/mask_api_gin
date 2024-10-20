package service

import (
	"errors"
	"fmt"
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/framework/utils/parse"
	systemModel "mask_api_gin/src/modules/system/model"
	systemService "mask_api_gin/src/modules/system/service"
)

// 实例化服务层 RegisterImpl 结构体
var NewRegisterImpl = &RegisterImpl{
	sysUserService:   systemService.NewSysUserImpl,
	sysConfigService: systemService.NewSysConfigImpl,
	sysRoleService:   systemService.NewSysRoleImpl,
}

// 账号注册操作处理 服务层处理
type RegisterImpl struct {
	// 用户信息服务
	sysUserService systemService.ISysUser
	// 参数配置服务
	sysConfigService systemService.ISysConfig
	// 角色服务
	sysRoleService systemService.ISysRole
}

// ValidateCaptcha 校验验证码
func (s *RegisterImpl) ValidateCaptcha(code, uuid string) error {
	// 验证码检查，从数据库配置获取验证码开关 true开启，false关闭
	captchaEnabledStr := s.sysConfigService.SelectConfigValueByKey("sys.account.captchaEnabled")
	if !parse.Boolean(captchaEnabledStr) {
		return nil
	}
	if code == "" || uuid == "" {
		return errors.New("验证码信息错误")
	}
	verifyKey := cachekey.CAPTCHA_CODE_KEY + uuid
	captcha, err := redis.Get("", verifyKey)
	if captcha == "" || err != nil {
		return errors.New("验证码已失效")
	}
	redis.Del("", verifyKey)
	if captcha != code {
		return errors.New("验证码错误")
	}
	return nil
}

// ByUserName 账号注册
func (s *RegisterImpl) ByUserName(username, password, userType string) (string, error) {
	// 是否开启用户注册功能 true开启，false关闭
	registerUserStr := s.sysConfigService.SelectConfigValueByKey("sys.account.registerUser")
	captchaEnabled := parse.Boolean(registerUserStr)
	if !captchaEnabled {
		return "", fmt.Errorf("注册用户【%s】失败，很抱歉，系统已关闭外部用户注册通道", username)
	}

	// 检查用户登录账号是否唯一
	uniqueUserName := s.sysUserService.CheckUniqueUserName(username, "")
	if !uniqueUserName {
		return "", fmt.Errorf("注册用户【%s】失败，注册账号已存在", username)
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
		return insertId, nil
	}
	return "", fmt.Errorf("注册用户【%s】失败，请联系系统管理人员", username)
}

// registerRoleInit 注册初始角色
func (s *RegisterImpl) registerRoleInit(userType string) []string {
	if userType == "sys" {
		return []string{}
	}
	return []string{}
}

// registerPostInit 注册初始岗位
func (s *RegisterImpl) registerPostInit(userType string) []string {
	if userType == "sys" {
		return []string{}
	}
	return []string{}
}
