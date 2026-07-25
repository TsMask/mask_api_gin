package repository

import (
	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysRoleMenu 实例化数据层
var NewSysRoleMenu = &SysRoleMenu{}

// SysRoleMenu 角色与菜单关联表 数据层处理
type SysRoleMenu struct{}

// ExistRoleByMenuId 存在角色使用数量By菜单
func (r SysRoleMenu) ExistRoleByMenuId(menuId int64) (int64, error) {
	if menuId <= 0 {
		return 0, nil
	}
	tx := database.DB("").Model(&model.SysRoleMenu{})
	tx = tx.Where("menu_id = ?", menuId)
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

// DeleteByRoleIds 批量删除关联By角色
func (r SysRoleMenu) DeleteByRoleIds(d *gorm.DB, roleIds []int64) (int64, error) {
	if len(roleIds) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.Where("role_id in ?", roleIds)
	// 执行删除
	if err := tx.Delete(&model.SysRoleMenu{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteByMenuIds 批量删除关联By菜单
func (r SysRoleMenu) DeleteByMenuIds(d *gorm.DB, menuIds []int64) (int64, error) {
	if len(menuIds) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.Where("menu_id in ?", menuIds)
	// 执行删除
	if err := tx.Delete(&model.SysRoleMenu{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// BatchInsert 批量新增信息
func (r SysRoleMenu) BatchInsert(d *gorm.DB, sysRoleMenus []model.SysRoleMenu) (int64, error) {
	if len(sysRoleMenus) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.CreateInBatches(sysRoleMenus, 500)
	if err := tx.Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}
