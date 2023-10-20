package service

import "mask_api_gin/src/modules/system/model"

// ISysLogLogin 系统登录日志 服务层接口
type ISysLogLogin interface {
	// SelectSysLogLoginPage 分页查询系统登录日志集合
	SelectSysLogLoginPage(query map[string]any) map[string]any

	// SelectSysLogLoginList 查询系统登录日志集合
	SelectSysLogLoginList(sysLogLogin model.SysLogLogin) []model.SysLogLogin

	// InsertSysLogLogin 新增系统登录日志
	InsertSysLogLogin(sysLogLogin model.SysLogLogin) string

	// DeleteSysLogLoginByIds 批量删除系统登录日志
	DeleteSysLogLoginByIds(loginIds []string) int64

	// CleanSysLogLogin 清空系统登录日志
	CleanSysLogLogin() error

	// CreateSysLogLogin 创建系统登录日志
	CreateSysLogLogin(userName, status, msg string, ilobArgs ...string) string
}
