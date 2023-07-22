package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
)

// SysRoleMenuImpl 角色与菜单关联表 数据层处理
var SysRoleMenuImpl = &sysRoleMenuImpl{
	selectSql: "",
}

type sysRoleMenuImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// CheckMenuExistRole 查询菜单分配给角色使用数量
func (r *sysRoleMenuImpl) CheckMenuExistRole(menuId string) int64 {
	querySql := "select count(1) as 'total' from sys_role_menu where menu_id = ?"
	results, err := datasource.RawDB("", querySql, []interface{}{menuId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return results[0]["total"].(int64)
	}
	return 0
}

// DeleteRoleMenuByRoleId 通过角色ID删除角色和菜单关联
func (r *sysRoleMenuImpl) DeleteRoleMenuByRoleId(roleId string) int {
	return 0
}

// DeleteRoleMenu 批量删除角色菜单关联信息
func (r *sysRoleMenuImpl) DeleteRoleMenu(roleIds []string) int {
	return 0
}

// BatchRoleMenu 批量新增角色菜单信息
func (r *sysRoleMenuImpl) BatchRoleMenu(sysRoleMenus []model.SysRoleMenu) int {
	return 0
}
