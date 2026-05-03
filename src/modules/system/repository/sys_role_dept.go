package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysRoleDept 实例化数据层
var NewSysRoleDept = &SysRoleDept{}

// SysRoleDept 角色与部门关联表 数据层处理
type SysRoleDept struct{}

// DeleteByRoleIds 批量删除信息By角色
func (r SysRoleDept) DeleteByRoleIds(d *gorm.DB, roleIds []int64) int64 {
	if len(roleIds) <= 0 {
		return 0
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Where("role_id in ?", roleIds)
	// 执行删除
	if err := tx.Delete(&model.SysRoleDept{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByDeptIds 批量删除信息By部门
func (r SysRoleDept) DeleteByDeptIds(d *gorm.DB, deptIds []int64) int64 {
	if len(deptIds) <= 0 {
		return 0
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Where("dept_id in ?", deptIds)
	// 执行删除
	if err := tx.Delete(&model.SysRoleDept{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// BatchInsert 批量新增信息
func (r SysRoleDept) BatchInsert(d *gorm.DB, sysRoleDepts []model.SysRoleDept) int64 {
	if len(sysRoleDepts) <= 0 {
		return 0
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.CreateInBatches(sysRoleDepts, 500)
	if err := tx.Error; err != nil {
		logger.Errorf("insert batch err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
