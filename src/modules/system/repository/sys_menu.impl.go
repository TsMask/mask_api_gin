package repository

import (
	"mask_api_gin/src/framework/constants/menu"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	repoUtils "mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strconv"
	"strings"
)

// SysMenuImpl 菜单表 数据层处理
var SysMenuImpl = &sysMenuImpl{
	selectSql: `select 
	menu_id, menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, ifnull(perms,'') as perms, icon, menu_sort, create_time, remark
	from sys_menu`,
	resultMap: map[string]string{
		"menu_id":     "MenuID",
		"menu_name":   "MenuName",
		"parent_name": "ParentName",
		"parent_id":   "ParentID",
		"path":        "Path",
		"menu_sort":   "MenuSort",
		"component":   "Component",
		"is_frame":    "IsFrame",
		"is_cache":    "IsCache",
		"menu_type":   "MenuType",
		"visible":     "Visible",
		"status":      "Status",
		"perms":       "Perms",
		"icon":        "Icon",
		"create_by":   "CreateBy",
		"create_time": "CreateTime",
		"update_by":   "UpdateBy",
		"update_time": "UpdateTime",
		"remark":      "Remark",
	},
}

type sysMenuImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *sysMenuImpl) convertResultRows(rows []map[string]interface{}) []model.SysMenu {
	arr := make([]model.SysMenu, 0)
	for _, row := range rows {
		sysMenu := model.SysMenu{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repoUtils.SetFieldValue(&sysMenu, keyMapper, value)
			}
		}
		arr = append(arr, sysMenu)
	}
	return arr
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
	var params []interface{}
	var querySql string

	if userId == "*" {
		// 管理员全部菜单
		querySql = `select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.menu_sort, m.create_time, m.remark
		from sys_menu m where m.menu_type in (?,?) and m.status = '1'
		order by m.parent_id, m.menu_sort`
		params = append(params, menu.TYPE_DIR)
		params = append(params, menu.TYPE_MENU)
	} else {
		// 用户ID权限
		querySql = `select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.menu_sort, m.create_time, m.remark
		from sys_menu m
		left join sys_role_menu rm on m.menu_id = rm.menu_id
		left join sys_user_role ur on rm.role_id = ur.role_id
		left join sys_role ro on ur.role_id = ro.role_id
		left join sys_user u on ur.user_id = u.user_id
		where u.user_id = ? and m.menu_type in (?,?) and m.status = '1'  AND ro.status = 0
		order by m.parent_id, m.menu_sort`
		params = append(params, userId)
		params = append(params, menu.TYPE_DIR)
		params = append(params, menu.TYPE_MENU)
	}

	// 查询结果
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysMenu{}
	}

	return r.convertResultRows(results)
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
	if whereSql == "" {
		return ""
	}

	// 查询数据
	querySql := "select menu_id as 'str' from sys_menu " + whereSql + " limit 1"
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
	}
	if len(results) > 0 {
		return strconv.FormatInt(results[0]["str"].(int64), 10)
	}
	return ""
}
