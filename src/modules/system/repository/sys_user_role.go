package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
)

// NewSysUserRole 实例化数据层
var NewSysUserRole = &SysUserRole{}

// SysUserRole 用户与角色关联表 数据层处理
type SysUserRole struct{}

// ExistUserByRoleId 存在用户使用数量
func (r SysUserRole) ExistUserByRoleId(roleId string) int64 {
	if roleId == "" {
		return 0
	}
	tx := db.DB("").Model(&model.SysUserRole{})
	tx = tx.Where("role_id = ?", roleId)
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return count
	}
	return count
}

// DeleteByUserIds 批量删除关联By用户
func (r SysUserRole) DeleteByUserIds(userIds []string) int64 {
	if len(userIds) <= 0 {
		return 0
	}
	tx := db.DB("").Where("user_id in ?", userIds)
	// 执行删除
	if err := tx.Delete(&model.SysUserRole{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByRoleId 批量删除关联By角色
func (r SysUserRole) DeleteByRoleId(roleId string, userIds []string) int64 {
	if roleId == "" || len(userIds) <= 0 {
		return 0
	}
	tx := db.DB("").Where("role_id = ?", roleId).Where("user_id in ?", userIds)
	// 执行删除
	if err := tx.Delete(&model.SysUserRole{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// BatchInsert 批量新增信息
func (r SysUserRole) BatchInsert(userRoles []model.SysUserRole) int64 {
	if len(userRoles) <= 0 {
		return 0
	}
	// 执行批量删除
	tx := db.DB("").CreateInBatches(userRoles, 500)
	if err := tx.Error; err != nil {
		logger.Errorf("delete batch err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
