package service

import "mask_api_gin/src/modules/system/model"

// ISysRoleDept 角色与部门关联表 数据层接口
type ISysRoleDept interface {
	// DeleteRoleDept 批量删除角色部门关联信息
	DeleteRoleDept(roleIds []string) int

	// BatchRoleDept 批量新增角色部门信息
	BatchRoleDept(sysRoleDepts []model.SysRoleDept) int
}
