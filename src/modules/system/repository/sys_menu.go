package repository

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"

	"time"
)

// NewSysMenu 实例化数据层
var NewSysMenu = &SysMenu{}

// SysMenu 菜单表 数据层处理
type SysMenu struct{}

// Select 查询集合 userId为0为系统管理员
func (r SysMenu) Select(sysMenu model.SysMenu, userId string) []model.SysMenu {
	tx := db.DB("").Model(&model.SysMenu{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if sysMenu.MenuName != "" {
		tx = tx.Where("menu_name like concat(?, '%')", sysMenu.MenuName)
	}
	if sysMenu.VisibleFlag != "" {
		tx = tx.Where("visible_flag = ?", sysMenu.VisibleFlag)
	}
	if sysMenu.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", sysMenu.StatusFlag)
	}

	// 个人菜单
	if userId != "0" {
		tx = tx.Where(`menu_id in (
		select menu_id from sys_role_menu where role_id in (
		select role_id from sys_user_role where user_id = ?
		))`, userId)
	}

	// 查询数据
	rows := []model.SysMenu{}
	if err := tx.Order("parent_id, menu_sort").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r SysMenu) SelectByIds(menuIds []string) []model.SysMenu {
	rows := []model.SysMenu{}
	if len(menuIds) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysMenu{})
	// 构建查询条件
	tx = tx.Where("menu_id in ? and del_flag = '0'", menuIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息
func (r SysMenu) Insert(sysMenu model.SysMenu) string {
	sysMenu.DelFlag = "0"
	if sysMenu.MenuId != "" {
		return ""
	}
	if sysMenu.Icon == "" {
		sysMenu.Icon = "#"
	}
	if sysMenu.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysMenu.UpdateBy = sysMenu.CreateBy
		sysMenu.UpdateTime = ms
		sysMenu.CreateTime = ms
	}

	// 根据菜单类型重置参数
	if sysMenu.MenuType == constants.MENU_TYPE_BUTTON {
		sysMenu.Component = ""
		sysMenu.Perms = ""
		sysMenu.Icon = "#"
		sysMenu.FrameFlag = "1"
		sysMenu.CacheFlag = "1"
		sysMenu.VisibleFlag = "1"
		sysMenu.StatusFlag = "1"
	} else if sysMenu.MenuType == constants.MENU_TYPE_DIR {
		sysMenu.Component = ""
		sysMenu.Perms = ""
	}

	// 执行插入
	if err := db.DB("").Create(&sysMenu).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return ""
	}
	return sysMenu.MenuId
}

// Update 修改信息
func (r SysMenu) Update(sysMenu model.SysMenu) int64 {
	if sysMenu.MenuId == "" {
		return 0
	}
	if sysMenu.Icon == "" {
		sysMenu.Icon = "#"
	}
	if sysMenu.UpdateBy != "" {
		sysMenu.UpdateTime = time.Now().UnixMilli()
	}

	// 根据菜单类型重置参数
	if sysMenu.MenuType == constants.MENU_TYPE_BUTTON {
		sysMenu.Component = ""
		sysMenu.Perms = ""
		sysMenu.Icon = "#"
		sysMenu.FrameFlag = "1"
		sysMenu.CacheFlag = "1"
		sysMenu.VisibleFlag = "1"
		sysMenu.StatusFlag = "1"
	} else if sysMenu.MenuType == constants.MENU_TYPE_DIR {
		sysMenu.Component = ""
		sysMenu.Perms = ""
	}

	tx := db.DB("").Model(&model.SysMenu{})
	// 构建查询条件
	tx = tx.Where("menu_id = ?", sysMenu.MenuId)
	tx = tx.Omit("menu_id", "del_flag", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysMenu).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteById 删除信息 返回受影响行数
func (r SysMenu) DeleteById(menuId string) int64 {
	if menuId == "" {
		return 0
	}
	tx := db.DB("").Model(&model.SysMenu{})
	// 构建查询条件
	tx = tx.Where("menu_id = ?", menuId)
	// 执行更新删除标记
	if err := tx.Update("del_flag", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// ExistChildrenByMenuIdAndStatus 菜单下同状态存在子节点数量
func (r SysMenu) ExistChildrenByMenuIdAndStatus(menuId string, statusFlag string) int64 {
	if menuId == "" {
		return 0
	}
	tx := db.DB("").Model(&model.SysMenu{})
	// 构建查询条件
	tx = tx.Where("parent_id = ? and del_flag = '0'", menuId)
	if statusFlag != "" {
		tx = tx.Where("status_flag = ?", statusFlag)
		tx = tx.Where("menu_type in ?", []string{constants.MENU_TYPE_DIR, constants.MENU_TYPE_MENU})
	}
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return count
	}
	return count
}

// CheckUnique 检查信息是否唯一
func (r SysMenu) CheckUnique(sysMenu model.SysMenu) string {
	tx := db.DB("").Model(&model.SysMenu{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if sysMenu.ParentId != "" {
		tx = tx.Where("parent_id = ?", sysMenu.ParentId)
	}
	if sysMenu.MenuName != "" {
		tx = tx.Where("menu_name = ?", sysMenu.MenuName)
	}
	if sysMenu.MenuPath != "" {
		tx = tx.Where("menu_path = ?", sysMenu.MenuPath)
	}

	// 查询数据
	var id string = ""
	if err := tx.Select("menu_id").Limit(1).Find(&id).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return id
	}
	return id
}

// SelectPermsByUserId 根据用户ID查询权限标识
func (r SysMenu) SelectPermsByUserId(userId string) []string {
	rows := []string{}
	if userId == "" {
		return rows
	}
	tx := db.DB("").Table("sys_menu m")
	// 构建查询条件
	tx = tx.Distinct("m.perms").
		Joins("left join sys_role_menu rm on m.menu_id = rm.menu_id").
		Joins("left join sys_user_role ur on rm.role_id = ur.role_id").
		Joins("left join sys_role r on r.role_id = ur.role_id").
		Where("m.status_flag = '1' and m.perms != '' and r.status_flag = '1' and r.del_flag = '0'").
		Where("ur.user_id = ?", userId)

	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByRoleId 根据角色ID查询菜单树信息 TODO
func (r SysMenu) SelectByRoleId(roleId string, menuCheckStrictly bool) []string {
	if roleId == "" {
		return []string{}
	}

	tx := db.DB("").Model(&model.SysMenu{})
	tx = tx.Where("del_flag = '0'")
	tx = tx.Where("menu_id in (select menu_id from sys_role_menu where role_id = ?)", roleId)
	// 展开
	if menuCheckStrictly {
		tx = tx.Where(`menu_id not in (
		select m.parent_id from sys_menu m 
		inner join sys_role_menu rm on m.menu_id = rm.menu_id 
		and rm.role_id = ?
		)`, roleId)
	}

	// 查询数据
	rows := []string{}
	if err := tx.Distinct("menu_id").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectTreeByUserId 根据用户ID查询菜单 0为管理员查询全部菜单，其他为用户ID查询权限
func (r SysMenu) SelectTreeByUserId(userId string) []model.SysMenu {
	if userId == "" {
		return []model.SysMenu{}
	}

	tx := db.DB("").Model(&model.SysMenu{})
	tx = tx.Where("del_flag = '0'")
	// 管理员全部菜单
	if userId == "0" {
		tx = tx.Where("menu_type in ? and status_flag = '1'", []string{constants.MENU_TYPE_DIR, constants.MENU_TYPE_MENU})
	} else {
		// 用户ID权限
		tx = tx.Where(`menu_type in ? and status_flag = '1' 
		and menu_id in (
		select menu_id from sys_role_menu where role_id in (
		select role_id from sys_user_role where user_id = ?
		))`, []string{constants.MENU_TYPE_DIR, constants.MENU_TYPE_MENU}, userId)
	}

	// 查询数据
	rows := []model.SysMenu{}
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}
