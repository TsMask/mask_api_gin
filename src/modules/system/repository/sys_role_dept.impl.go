package repository

import "mask_api_gin/src/modules/system/model"

// SysRoleDeptImpl 角色与部门关联表 数据层处理
var SysRoleDeptImpl = &sysRoleDeptImpl{
	selectSql: "",
}

type sysRoleDeptImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// DeleteRoleDept 批量删除角色部门关联信息
func (r *sysRoleDeptImpl) DeleteRoleDept(roleIds []string) int {
	return 0
}

// BatchRoleDept 批量新增角色部门信息
func (r *sysRoleDeptImpl) BatchRoleDept(sysRoleDepts []model.SysRoleDept) int {
	return 0
}
