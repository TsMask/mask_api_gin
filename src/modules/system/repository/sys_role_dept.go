package repository

import (
	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysRoleDept 实例化数据层
var NewSysRoleDept = &SysRoleDept{}

// SysRoleDept 角色与部门关联表 数据层处理
type SysRoleDept struct{}

// DeleteByRoleIds 批量删除信息By角色
func (r SysRoleDept) DeleteByRoleIds(d *gorm.DB, roleIds []int64) (int64, error) {
	if len(roleIds) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.Where("role_id in ?", roleIds)
	// 执行删除
	if err := tx.Delete(&model.SysRoleDept{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteByDeptIds 批量删除信息By部门
func (r SysRoleDept) DeleteByDeptIds(d *gorm.DB, deptIds []int64) (int64, error) {
	if len(deptIds) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.Where("dept_id in ?", deptIds)
	// 执行删除
	if err := tx.Delete(&model.SysRoleDept{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// BatchInsert 批量新增信息
func (r SysRoleDept) BatchInsert(d *gorm.DB, sysRoleDepts []model.SysRoleDept) (int64, error) {
	if len(sysRoleDepts) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.CreateInBatches(sysRoleDepts, 500)
	if err := tx.Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}
