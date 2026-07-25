package repository

import (
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysMenu 实例化数据层
var NewSysMenu = &SysMenu{}

// SysMenu 菜单表 数据层处理
type SysMenu struct{}

// Select 查询集合 userId为0是系统管理员
func (r SysMenu) Select(sysMenu model.SysMenu, userId int64) ([]model.SysMenu, error) {
	tx := database.DB("").Model(&model.SysMenu{})
	tx = tx.Where("menu_del = '0'")
	// 查询条件拼接
	if sysMenu.MenuName != "" {
		tx = tx.Where("menu_name like ?", sysMenu.MenuName+"%")
	}
	if sysMenu.MenuVisible != "" {
		tx = tx.Where("menu_visible = ?", sysMenu.MenuVisible)
	}
	if sysMenu.MenuStatus != "" {
		tx = tx.Where("menu_status = ?", sysMenu.MenuStatus)
	}

	// 个人菜单（sys_menu 主键 id；sys_role_menu/sys_user_role 外键不变）
	if userId > 0 {
		tx = tx.Where(`id in (
		select menu_id from sys_role_menu where role_id in (
		select role_id from sys_user_role where user_id = ?
		))`, userId)
	}

	// 查询数据
	rows := []model.SysMenu{}
	tx.Order("parent_id, menu_sort")
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectByIds 通过ID查询信息
func (r SysMenu) SelectByIds(menuIds []int64) ([]model.SysMenu, error) {
	rows := []model.SysMenu{}
	if len(menuIds) <= 0 {
		return rows, nil
	}
	tx := database.DB("").Model(&model.SysMenu{})
	// 构建查询条件
	tx = tx.Where("id in ? and menu_del = '0'", menuIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// Insert 新增信息
func (r SysMenu) Insert(sysMenu model.SysMenu) (int64, error) {
	sysMenu.MenuDel = "0"
	if sysMenu.ID > 0 {
		return 0, nil
	}

	// 执行插入
	if err := database.DB("").Create(&sysMenu).Error; err != nil {
		return 0, err
	}
	return sysMenu.ID, nil
}

// Update 修改信息
func (r SysMenu) Update(sysMenu model.SysMenu) (int64, error) {
	if sysMenu.ID <= 0 {
		return 0, nil
	}

	tx := database.DB("").Model(&model.SysMenu{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysMenu.ID)
	tx = tx.Omit("id", "create_by", "create_time", "menu_del")
	// 执行更新
	if err := tx.Updates(sysMenu).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteById 删除信息 返回受影响行数（软删除）
func (r SysMenu) DeleteById(d *gorm.DB, menuId int64) (int64, error) {
	if menuId <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.Model(&model.SysMenu{})
	// 构建查询条件
	tx = tx.Where("id = ?", menuId)
	// 执行更新删除标记
	if err := tx.Update("menu_del", "1").Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// CheckUnique 检查信息是否唯一
func (r SysMenu) CheckUnique(sysMenu model.SysMenu) (int64, error) {
	tx := database.DB("").Model(&model.SysMenu{})
	tx = tx.Where("menu_del = '0'")
	// 查询条件拼接
	if sysMenu.ParentId >= 0 {
		tx = tx.Where("parent_id = ?", sysMenu.ParentId)
	}
	if sysMenu.MenuName != "" {
		tx = tx.Where("menu_name = ?", sysMenu.MenuName)
	}
	if sysMenu.MenuPath != "" {
		tx = tx.Where("menu_path = ?", sysMenu.MenuPath)
	}

	// 查询数据
	var id int64 = 0
	if err := tx.Select("id").Limit(1).Find(&id).Error; err != nil {
		return id, err
	}
	return id, nil
}

// ExistChildrenByMenuIdAndStatus 菜单下同状态存在子节点数量
func (r SysMenu) ExistChildrenByMenuIdAndStatus(menuId int64, menuStatus string) (int64, error) {
	if menuId <= 0 {
		return 0, nil
	}
	tx := database.DB("").Model(&model.SysMenu{})
	// 构建查询条件
	tx = tx.Where("parent_id = ? and menu_del = '0'", menuId)
	if menuStatus != "" {
		tx = tx.Where("menu_status = ?", menuStatus)
		tx = tx.Where("menu_type in ?", []string{constants.MENU_TYPE_DIR, constants.MENU_TYPE_MENU})
	}
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

// SelectPermsByUserId 根据用户ID查询权限标识
func (r SysMenu) SelectPermsByUserId(userId int64) ([]string, error) {
	rows := []string{}
	if userId <= 0 {
		return rows, nil
	}

	tx := database.DB("").Model(&model.SysMenu{}).
		Table("sys_menu m").
		Distinct("m.menu_perms").
		Joins("LEFT JOIN sys_role_menu rm ON m.id = rm.menu_id").
		Joins("LEFT JOIN sys_user_role ur ON rm.role_id = ur.role_id").
		Joins("LEFT JOIN sys_role r ON r.id = ur.role_id").
		Where("m.menu_status = '1' and m.menu_perms != ''").
		Where("r.role_status = '1'").
		Where("ur.user_id = ?", userId)

	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectByRoleId 根据角色ID查询菜单树信息
func (r SysMenu) SelectByRoleId(roleId int64, menuCheckStrictly bool) ([]int64, error) {
	if roleId <= 0 {
		return []int64{}, nil
	}

	tx := database.DB("").Model(&model.SysMenu{})
	tx = tx.Where("menu_del = '0'")
	tx = tx.Where("id in (select menu_id from sys_role_menu where role_id = ?)", roleId)
	// 父子互相关联显示，取所有子节点
	if menuCheckStrictly {
		tx = tx.Where(`id not in (
		select m.parent_id from sys_menu m
		inner join sys_role_menu rm on m.id = rm.menu_id
		and rm.role_id = ?
		)`, roleId)
	}

	// 查询数据
	rows := []int64{}
	if err := tx.Distinct("id").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectTreeByUserId 根据用户ID查询菜单 0为管理员查询全部菜单，其他为用户ID查询权限
func (r SysMenu) SelectTreeByUserId(userId int64) ([]model.SysMenu, error) {
	if userId < 0 {
		return []model.SysMenu{}, nil
	}

	tx := database.DB("").Model(&model.SysMenu{})
	tx = tx.Where("menu_del = '0'")
	// 管理员全部菜单
	if userId == 0 {
		tx = tx.Where("menu_type in ? and menu_status = '1'", []string{constants.MENU_TYPE_DIR, constants.MENU_TYPE_MENU})
	} else {
		// 用户ID权限
		tx = tx.Where(`menu_type in ? and menu_status = '1'
		and id in (
		select menu_id from sys_role_menu where role_id in (
		select role_id from sys_user_role where user_id = ?
		))`, []string{constants.MENU_TYPE_DIR, constants.MENU_TYPE_MENU}, userId)
	}

	// 查询数据
	rows := []model.SysMenu{}
	tx.Order("parent_id, menu_sort")
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}
