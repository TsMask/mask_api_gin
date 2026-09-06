package service

import (
	"fmt"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/generate"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/shared/regular"
	systemModel "mask_api_gin/src/modules/system/model"
	systemService "mask_api_gin/src/modules/system/service"
)

// NewRegister 实例化服务层
var NewRegister = &Register{
	sysUserService:   systemService.NewSysUser,
	sysConfigService: systemService.NewSysConfig,
	accountService:   NewAccount,
}

// Register 账号注册操作 服务层处理
type Register struct {
	sysUserService   *systemService.SysUser   // 用户信息服务
	sysConfigService *systemService.SysConfig // 参数配置服务
	accountService   *Account                 // 账号身份操作服务
}

// ValidateCaptcha 校验账号注册图片验证码
func (s Register) ValidateCaptcha(code, uuid string) error {
	return s.accountService.ValidateCaptcha(code, uuid)
}

// ValidateSmsCaptcha 校验手机号注册短信验证码
func (s Register) ValidateSmsCaptcha(phone, code, uuid string) error {
	return s.accountService.ValidateCaptchaSMS(phone, code, uuid)
}

// ByUserName 账号注册
func (s Register) ByUserName(username, password string) (int64, error) {
	registerEnabled := parse.Boolean(s.sysConfigService.FindValueByKey("sys.user.register"))
	if !registerEnabled {
		return 0, fmt.Errorf("很抱歉，系统已关闭外部用户注册通道")
	}

	if !regular.ValidUsername(username) || regular.ValidMobile(username) {
		return 0, fmt.Errorf("用户账号格式不正确")
	}
	uniqueUserName := s.sysUserService.CheckUniqueByUserName(username, 0)
	if !uniqueUserName {
		return 0, fmt.Errorf("注册用户【%s】失败，注册账号已存在", username)
	}

	sysUser := systemModel.SysUser{
		UserName:     username,
		UserNick:     username,
		UserPassword: password,
		UserGender:   "0",
		UserStatus:   constants.STATUS_YES,
		DeptId:       100,
		CreateBy:     "reg-user",
		RoleIds:      s.registerRoleInit(),
		PostIds:      s.registerPostInit(),
	}

	insertId := s.sysUserService.Insert(sysUser)
	if insertId > 0 {
		return insertId, nil
	}
	return 0, fmt.Errorf("注册用户【%s】失败，请联系系统管理员", username)
}

// ByPhone 手机号注册
func (s Register) ByPhone(phone string) (int64, error) {
	registerEnabled := parse.Boolean(s.sysConfigService.FindValueByKey("sys.user.register"))
	if !registerEnabled {
		return 0, fmt.Errorf("很抱歉，系统已关闭外部用户注册通道")
	}

	if !regular.ValidMobile(phone) {
		return 0, fmt.Errorf("手机号格式不正确")
	}
	if !s.sysUserService.CheckUniqueByUserName(phone, 0) {
		return 0, fmt.Errorf("注册用户【%s】失败，注册账号已存在", phone)
	}
	if !s.sysUserService.CheckUniqueByPhone(phone, 0) {
		return 0, fmt.Errorf("注册用户【%s】失败，手机号码已存在", phone)
	}

	sysUser := systemModel.SysUser{
		UserName:     phone,
		UserPhone:    phone,
		UserNick:     phone,
		UserPassword: generate.String(32),
		UserGender:   "0",
		UserStatus:   constants.STATUS_YES,
		DeptId:       100,
		CreateBy:     "reg-phone",
		RoleIds:      s.registerRoleInit(),
		PostIds:      s.registerPostInit(),
	}

	insertId := s.sysUserService.Insert(sysUser)
	if insertId > 0 {
		return insertId, nil
	}
	return 0, fmt.Errorf("注册用户【%s】失败，请联系系统管理员", phone)
}

// registerRoleInit 注册初始化角色
func (s Register) registerRoleInit() []int64 {
	return []int64{}
}

// registerPostInit 注册初始化岗位
func (s Register) registerPostInit() []int64 {
	return []int64{}
}
