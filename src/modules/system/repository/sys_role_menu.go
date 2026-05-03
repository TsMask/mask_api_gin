package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysRoleMenu 实例化数据层
var NewSysRoleMenu = &SysRoleMenu{}

// SysRoleMenu 角色与菜单关联表 数据层处理
type SysRoleMenu struct{}

// ExistRoleByMenuId 存在角色使用数量By菜单
func (r SysRoleMenu) ExistRoleByMenuId(menuId int64) int64 {
	if menuId <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysRoleMenu{})
	tx = tx.Where("menu_id = ?", menuId)
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return count
	}
	return count
}

// DeleteByRoleIds 批量删除关联By角色
func (r SysRoleMenu) DeleteByRoleIds(d *gorm.DB, roleIds []int64) int64 {
	if len(roleIds) <= 0 {
		return 0
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Where("role_id in ?", roleIds)
	// 执行删除
	if err := tx.Delete(&model.SysRoleMenu{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByMenuIds 批量删除关联By菜单
func (r SysRoleMenu) DeleteByMenuIds(d *gorm.DB, menuIds []int64) int64 {
	if len(menuIds) <= 0 {
		return 0
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Where("menu_id in ?", menuIds)
	// 执行删除
	if err := tx.Delete(&model.SysRoleMenu{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// BatchInsert 批量新增信息
func (r SysRoleMenu) BatchInsert(d *gorm.DB, sysRoleMenus []model.SysRoleMenu) int64 {
	if len(sysRoleMenus) <= 0 {
		return 0
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.CreateInBatches(sysRoleMenus, 500)
	if err := tx.Error; err != nil {
		logger.Errorf("insert batch err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
