package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysUserRole 实例化数据层
var NewSysUserRole = &SysUserRole{}

// SysUserRole 用户与角色关联表 数据层处理
type SysUserRole struct{}

// ExistUserByRoleId 存在用户使用数量
func (r SysUserRole) ExistUserByRoleId(roleId int64) (int64, error) {
	if roleId <= 0 {
		return 0, nil
	}
	tx := db.DB("").Model(&model.SysUserRole{})
	tx = tx.Where("role_id = ?", roleId)
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

// SelectUserIdsByRoleIds 通过角色ID列表查询关联的用户ID
func (r SysUserRole) SelectUserIdsByRoleIds(roleIds []int64) ([]int64, error) {
	if len(roleIds) == 0 {
		return nil, nil
	}
	tx := db.DB("").Model(&model.SysUserRole{})
	tx = tx.Where("role_id in ?", roleIds)
	var userIds []int64
	if err := tx.Distinct("user_id").Pluck("user_id", &userIds).Error; err != nil {
		return nil, err
	}
	return userIds, nil
}

// DeleteByUserIds 批量删除关联By用户
func (r SysUserRole) DeleteByUserIds(d *gorm.DB, userIds []int64) (int64, error) {
	if len(userIds) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Where("user_id in ?", userIds)
	// 执行删除
	if err := tx.Delete(&model.SysUserRole{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteByRoleId 批量删除关联By角色
func (r SysUserRole) DeleteByRoleId(d *gorm.DB, roleId int64, userIds []int64) (int64, error) {
	if roleId <= 0 || len(userIds) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Where("role_id = ?", roleId).Where("user_id in ?", userIds)
	// 执行删除
	if err := tx.Delete(&model.SysUserRole{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// BatchInsert 批量新增信息
func (r SysUserRole) BatchInsert(d *gorm.DB, sysUserRoles []model.SysUserRole) (int64, error) {
	if len(sysUserRoles) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.CreateInBatches(sysUserRoles, 500)
	if err := tx.Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}
