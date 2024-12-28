package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
)

// NewSysRoleDept 实例化数据层
var NewSysRoleDept = &SysRoleDept{}

// SysRoleDept 角色与部门关联表 数据层处理
type SysRoleDept struct{}

// DeleteByRoleIds 批量删除信息By角色
func (r SysRoleDept) DeleteByRoleIds(roleIds []int64) int64 {
	if len(roleIds) <= 0 {
		return 0
	}
	tx := db.DB("").Where("role_id in ?", roleIds)
	// 执行删除
	if err := tx.Delete(&model.SysRoleDept{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByDeptIds 批量删除信息By部门
func (r SysRoleDept) DeleteByDeptIds(deptIds []int64) int64 {
	if len(deptIds) <= 0 {
		return 0
	}
	tx := db.DB("").Where("dept_id in ?", deptIds)
	// 执行删除
	if err := tx.Delete(&model.SysRoleDept{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// BatchInsert 批量新增信息
func (r SysRoleDept) BatchInsert(roleDepts []model.SysRoleDept) int64 {
	if len(roleDepts) <= 0 {
		return 0
	}
	// 执行批量删除
	tx := db.DB("").CreateInBatches(roleDepts, 500)
	if err := tx.Error; err != nil {
		logger.Errorf("delete batch err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
