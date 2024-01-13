package repository

import (
	"fmt"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// 实例化数据层 SysRoleMenuImpl 结构体
var NewSysRoleMenuImpl = &SysRoleMenuImpl{}

// SysRoleMenuImpl 角色与菜单关联表 数据层处理
type SysRoleMenuImpl struct{}

// CheckMenuExistRole 查询菜单分配给角色使用数量
func (r *SysRoleMenuImpl) CheckMenuExistRole(menuId string) int64 {
	querySql := "select count(1) as 'total' from sys_role_menu where menu_id = ?"
	results, err := datasource.RawDB("", querySql, []any{menuId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return parse.Number(results[0]["total"])
	}
	return 0
}

// DeleteRoleMenu 批量删除角色和菜单关联
func (r *SysRoleMenuImpl) DeleteRoleMenu(roleIds []string) int64 {
	placeholder := datasource.KeyPlaceholderByQuery(len(roleIds))
	sql := "delete from sys_role_menu where role_id in (" + placeholder + ")"
	parameters := datasource.ConvertIdsSlice(roleIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// DeleteMenuRole 批量删除菜单和角色关联
func (r *SysRoleMenuImpl) DeleteMenuRole(menuIds []string) int64 {
	placeholder := datasource.KeyPlaceholderByQuery(len(menuIds))
	sql := "delete from sys_role_menu where menu_id in (" + placeholder + ")"
	parameters := datasource.ConvertIdsSlice(menuIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// BatchRoleMenu 批量新增角色菜单信息
func (r *SysRoleMenuImpl) BatchRoleMenu(sysRoleMenus []model.SysRoleMenu) int64 {
	keyValues := make([]string, 0)
	for _, item := range sysRoleMenus {
		keyValues = append(keyValues, fmt.Sprintf("(%s,%s)", item.RoleID, item.MenuID))
	}
	sql := "insert into sys_role_menu(role_id, menu_id) values " + strings.Join(keyValues, ",")
	results, err := datasource.ExecDB("", sql, nil)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
