package repository

import "mask_api_gin/src/modules/system/model"

// ISysLogLoginRepository 系统登录日志表 数据层接口
type ISysLogLoginRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(query map[string]any) map[string]any

	// Select 查询集合
	Select(sysLogLogin model.SysLogLogin) []model.SysLogLogin

	// Insert 新增信息
	Insert(sysLogLogin model.SysLogLogin) string

	// DeleteByIds 批量删除信息
	DeleteByIds(loginIds []string) int64

	// Clean 清空信息
	Clean() error
}
