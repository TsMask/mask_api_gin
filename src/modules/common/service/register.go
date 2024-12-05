package service

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/utils/parse"
	systemModel "mask_api_gin/src/modules/system/model"
	systemService "mask_api_gin/src/modules/system/service"

	"fmt"
)

// NewRegister 实例化服务层
var NewRegister = &Register{
	sysUserService:   systemService.NewSysUser,
	sysConfigService: systemService.NewSysConfig,
	sysRoleService:   systemService.NewSysRole,
}

// Register 账号注册操作 服务层处理
type Register struct {
	sysUserService   *systemService.SysUser   // 用户信息服务
	sysConfigService *systemService.SysConfig // 参数配置服务
	sysRoleService   *systemService.SysRole   // 角色服务
}

// ValidateCaptcha 校验验证码
func (s Register) ValidateCaptcha(code, uuid string) error {
	// 验证码检查，从数据库配置获取验证码开关 true开启，false关闭
	captchaEnabledStr := s.sysConfigService.FindValueByKey("sys.account.captchaEnabled")
	if !parse.Boolean(captchaEnabledStr) {
		return nil
	}
	if code == "" || uuid == "" {
		return fmt.Errorf("captcha empty")
	}
	verifyKey := constants.CACHE_CAPTCHA_CODE + uuid
	captcha, _ := redis.Get("", verifyKey)
	if captcha == "" {
		return fmt.Errorf("captcha expire")
	}
	_ = redis.Del("", verifyKey)
	if captcha != code {
		return fmt.Errorf("captcha error")
	}
	return nil
}

// ByUserName 账号注册
func (s Register) ByUserName(username, password string) (string, error) {
	// 是否开启用户注册功能 true开启，false关闭
	registerUserStr := s.sysConfigService.FindValueByKey("sys.account.registerUser")
	captchaEnabled := parse.Boolean(registerUserStr)
	if !captchaEnabled {
		return "", fmt.Errorf("很抱歉，系统已关闭外部用户注册通道")
	}

	// 检查用户登录账号是否唯一
	uniqueUserName := s.sysUserService.CheckUniqueByUserName(username, "")
	if !uniqueUserName {
		return "", fmt.Errorf("注册用户【%s】失败，注册账号已存在", username)
	}

	sysUser := systemModel.SysUser{
		UserName:   username,
		NickName:   username,             // 昵称使用名称账号
		Password:   password,             // 原始密码
		Sex:        "0",                  // 性别未选择
		StatusFlag: constants.STATUS_YES, // 账号状态激活
		DeptId:     "100",                // 归属部门为根节点
		CreateBy:   "register",           // 创建来源
	}

	// 新增用户的角色管理
	sysUser.RoleIds = s.registerRoleInit()
	// 新增用户的岗位管理
	sysUser.PostIds = s.registerPostInit()

	insertId := s.sysUserService.Insert(sysUser)
	if insertId != "" {
		return insertId, nil
	}
	return "", fmt.Errorf("注册用户【%s】失败，请联系系统管理人员", username)
}

// registerRoleInit 注册初始角色
func (s Register) registerRoleInit() []string {
	return []string{}
}

// registerPostInit 注册初始岗位
func (s Register) registerPostInit() []string {
	return []string{}
}
