package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// SysMenuImpl 菜单表 数据层处理
var SysMenuImpl = &sysMenuImpl{
	selectSql: "",
}

type sysMenuImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// SelectMenuList 查询系统菜单列表
func (r *sysMenuImpl) SelectMenuList(sysMenu model.SysMenu, userId string) []model.SysMenu {
	return []model.SysMenu{}
}

// SelectMenuPermsByRoleId 根据角色ID查询权限
func (r *sysMenuImpl) SelectMenuPermsByRoleId(roleId string) []string {
	return []string{}
}

// SelectMenuPermsByUserId 根据用户ID查询权限
func (r *sysMenuImpl) SelectMenuPermsByUserId(userId string) []string {
	querySql := `select distinct m.perms as 'str' from sys_menu m 
    left join sys_role_menu rm on m.menu_id = rm.menu_id 
    left join sys_user_role ur on rm.role_id = ur.role_id 
    left join sys_role r on r.role_id = ur.role_id
	where m.status = '1' and r.status = '1' and ur.user_id = ? `

	// 查询结果
	results, err := datasource.RawDB("", querySql, []interface{}{userId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []string{}
	}

	// 读取结果
	rows := make([]string, 0)
	for _, m := range results {
		rows = append(rows, m["str"].(string))
	}
	return rows
}

// SelectMenuTreeByUserId 根据用户ID查询菜单
func (r *sysMenuImpl) SelectMenuTreeByUserId(userId string) []model.SysMenu {
	return []model.SysMenu{}
}

// SelectMenuListByRoleId 根据角色ID查询菜单树信息
func (r *sysMenuImpl) SelectMenuListByRoleId(roleId string, menuCheckStrictly bool) []string {
	return []string{}
}

// SelectMenuById 根据菜单ID查询信息
func (r *sysMenuImpl) SelectMenuById(menuId string) model.SysMenu {
	return model.SysMenu{}
}

// HasChildByMenuId 是否存在菜单子节点
func (r *sysMenuImpl) HasChildByMenuId(menuId string) int {
	return 0
}

// CheckMenuExistRole 查询菜单是否存在角色
func (r *sysMenuImpl) CheckMenuExistRole(menuId string) int {
	return 0
}

// InsertMenu 新增菜单信息
func (r *sysMenuImpl) InsertMenu(sysMenu model.SysMenu) string {
	return ""
}

// UpdateMenu 修改菜单信息
func (r *sysMenuImpl) UpdateMenu(sysMenu model.SysMenu) int {
	return 0
}

// DeleteMenuById 删除菜单管理信息
func (r *sysMenuImpl) DeleteMenuById(menuId string) int {
	return 0
}

// CheckUniqueMenu 校验菜单是否唯一
func (r *sysMenuImpl) CheckUniqueMenu(sysMenu model.SysMenu) string {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if sysMenu.MenuName != "" {
		conditions = append(conditions, "menu_name = ?")
		params = append(params, sysMenu.MenuName)
	}
	if sysMenu.ParentID != "" {
		conditions = append(conditions, "parent_id = ?")
		params = append(params, sysMenu.ParentID)
	}
	if sysMenu.Path != "" {
		conditions = append(conditions, "path = ?")
		params = append(params, sysMenu.Path)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数据
	querySql := "select menu_id as 'str' from sys_menu" + whereSql + "limit 1"
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
	}
	return results[0]["str"].(string)
}
