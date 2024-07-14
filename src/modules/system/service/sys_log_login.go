package service

import "mask_api_gin/src/modules/system/model"

// ISysLogLoginService 系统登录日志 服务层接口
type ISysLogLoginService interface {
	// FindByPage 分页查询列表数据
	FindByPage(query map[string]any) map[string]any

	// Find 查询数据
	Find(sysLogLogin model.SysLogLogin) []model.SysLogLogin

	// Insert 新增信息
	Insert(userName, status, msg string, ilobArr [4]string) string

	// DeleteByIds 批量删除信息
	DeleteByIds(loginIds []string) int64

	// Clean 清空系统登录日志
	Clean() error
}
