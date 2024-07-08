package repository

import "mask_api_gin/src/modules/system/model"

// ISysRoleDeptRepository 角色与部门关联表 数据层接口
type ISysRoleDeptRepository interface {
	// DeleteByRoleIds 批量删除关联By角色
	DeleteByRoleIds(roleIds []string) int64

	// DeleteByDeptIds 批量删除关联By部门
	DeleteByDeptIds(deptIds []string) int64

	// BatchInsert 批量新增信息
	BatchInsert(arr []model.SysRoleDept) int64
}
