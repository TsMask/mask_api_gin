package repository

import (
	"fmt"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// NewSysRoleMenu 实例化数据层
var NewSysRoleMenu = &SysRoleMenu{}

// SysRoleMenu 角色与菜单关联表 数据层处理
type SysRoleMenu struct{}

// ExistRoleByMenuId 存在角色使用数量
func (r SysRoleMenu) ExistRoleByMenuId(menuId string) int64 {
	querySql := "select count(1) as 'total' from sys_role_menu where menu_id = ?"
	results, err := db.RawDB("", querySql, []any{menuId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	return parse.Number(results[0]["total"])
}

// DeleteByRoleIds 批量删除关联By角色
func (r SysRoleMenu) DeleteByRoleIds(roleIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(roleIds))
	sql := fmt.Sprintf("delete from sys_role_menu where role_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(roleIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// DeleteByMenuIds 批量删除关联By菜单
func (r SysRoleMenu) DeleteByMenuIds(menuIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(menuIds))
	sql := fmt.Sprintf("delete from sys_role_menu where menu_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(menuIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// BatchInsert 批量新增信息
func (r SysRoleMenu) BatchInsert(arr []model.SysRoleMenu) int64 {
	rm := make([]string, 0)
	for _, item := range arr {
		rm = append(rm, fmt.Sprintf("(%s,%s)", item.RoleId, item.MenuId))
	}
	sql := fmt.Sprintf("insert into sys_role_menu(role_id, menu_id) values %s", strings.Join(rm, ","))
	results, err := db.ExecDB("", sql, nil)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
