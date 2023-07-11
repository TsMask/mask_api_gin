package repository

import "mask_api_gin/src/modules/monitor/model"

// SysLogininforImpl 系统登录访问表 数据层处理
var SysLogininforImpl = &sysLogininforImpl{
	selectSql: "",
}

type sysLogininforImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// SelectLogininforPage 分页查询系统登录日志集合
func (r *sysLogininforImpl) SelectLogininforPage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectLogininforList 查询系统登录日志集合
func (r *sysLogininforImpl) SelectLogininforList(sysLogininfor model.SysLogininfor) []model.SysLogininfor {
	return []model.SysLogininfor{}
}

// InsertLogininfor 新增系统登录日志
func (r *sysLogininforImpl) InsertLogininfor(sysLogininfor model.SysLogininfor) string {
	return r.selectSql
}

// DeleteLogininforByIds 批量删除系统登录日志
func (r *sysLogininforImpl) DeleteLogininforByIds(infoIds []string) int64 {
	return 0
}

// CleanLogininfor 清空系统登录日志
func (r *sysLogininforImpl) CleanLogininfor() error {
	return nil
}
