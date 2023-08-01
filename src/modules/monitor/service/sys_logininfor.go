package service

import "mask_api_gin/src/modules/monitor/model"

// ISysLogininfor 系统登录访问 服务层接口
type ISysLogininfor interface {
	// SelectLogininforPage 分页查询系统登录日志集合
	SelectLogininforPage(query map[string]any) map[string]any

	// SelectLogininforList 查询系统登录日志集合
	SelectLogininforList(sysLogininfor model.SysLogininfor) []model.SysLogininfor

	// InsertLogininfor 新增系统登录日志
	InsertLogininfor(sysLogininfor model.SysLogininfor) string

	// DeleteLogininforByIds 批量删除系统登录日志
	DeleteLogininforByIds(infoIds []string) int64

	// CleanLogininfor 清空系统登录日志
	CleanLogininfor() error

	// NewLogininfor 生成系统登录日志
	NewLogininfor(userName, status, msg string, ilobArgs ...string) string
}
