package repository

import (
	"fmt"
	"mask_api_gin/src/framework/constants/menu"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"

	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// 实例化数据层 SysMenuImpl 结构体
var NewSysMenuImpl = &SysMenuImpl{
	selectSql: `select 
	m.menu_id, m.menu_name, m.parent_id, m.menu_sort, m.path, m.component, m.is_frame, m.is_cache, m.menu_type, m.visible, m.status, ifnull(m.perms,'') as perms, m.icon, m.create_time, m.remark 
	from sys_menu m`,

	selectSqlByUser: `select distinct 
	m.menu_id, m.menu_name, m.parent_id, m.menu_sort, m.path, m.component, m.is_frame, m.is_cache, m.menu_type, m.visible, m.status, ifnull(m.perms,'') as perms, m.icon, m.create_time, m.remark
	from sys_menu m
	left join sys_role_menu rm on m.menu_id = rm.menu_id
	left join sys_user_role ur on rm.role_id = ur.role_id
	left join sys_role ro on ur.role_id = ro.role_id`,

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

// SysMenuImpl 菜单表 数据层处理
type SysMenuImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 查询视图用户对象SQL
	selectSqlByUser string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysMenuImpl) convertResultRows(rows []map[string]any) []model.SysMenu {
	arr := make([]model.SysMenu, 0)
	for _, row := range rows {
		sysMenu := model.SysMenu{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				datasource.SetFieldValue(&sysMenu, keyMapper, value)
			}
		}
		arr = append(arr, sysMenu)
	}
	return arr
}

// SelectMenuList 查询系统菜单列表
func (r *SysMenuImpl) SelectMenuList(sysMenu model.SysMenu, userId string) []model.SysMenu {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysMenu.MenuName != "" {
		conditions = append(conditions, "m.menu_name like concat(?, '%')")
		params = append(params, sysMenu.MenuName)
	}
	if sysMenu.Visible != "" {
		conditions = append(conditions, "m.visible = ?")
		params = append(params, sysMenu.Visible)
	}
	if sysMenu.Status != "" {
		conditions = append(conditions, "m.status = ?")
		params = append(params, sysMenu.Status)
	}

	fromSql := r.selectSql

	// 个人菜单
	if userId != "*" {
		fromSql = r.selectSqlByUser
		conditions = append(conditions, "ur.user_id = ?")
		params = append(params, userId)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数据
	orderSql := " order by m.parent_id, m.menu_sort"
	querySql := fromSql + whereSql + orderSql
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysMenu{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectMenuPermsByUserId 根据用户ID查询权限
func (r *SysMenuImpl) SelectMenuPermsByUserId(userId string) []string {
	querySql := `select distinct m.perms as 'str' from sys_menu m 
    left join sys_role_menu rm on m.menu_id = rm.menu_id 
    left join sys_user_role ur on rm.role_id = ur.role_id 
    left join sys_role r on r.role_id = ur.role_id
	where m.status = '1' and m.perms != '' and r.status = '1' and ur.user_id = ? `

	// 查询结果
	results, err := datasource.RawDB("", querySql, []any{userId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []string{}
	}

	// 读取结果
	rows := make([]string, 0)
	for _, m := range results {
		rows = append(rows, fmt.Sprintf("%v", m["str"]))
	}
	return rows
}

// SelectMenuTreeByUserId 根据用户ID查询菜单
func (r *SysMenuImpl) SelectMenuTreeByUserId(userId string) []model.SysMenu {
	var params []any
	var querySql string

	if userId == "*" {
		// 管理员全部菜单
		querySql = r.selectSql + ` where 
		m.menu_type in (?,?) and m.status = '1'
		order by m.parent_id, m.menu_sort`
		params = append(params, menu.TYPE_DIR)
		params = append(params, menu.TYPE_MENU)
	} else {
		// 用户ID权限
		querySql = r.selectSqlByUser + ` where 
		m.menu_type in (?, ?) and m.status = '1'
		and ur.user_id = ? and ro.status = '1'
		order by m.parent_id, m.menu_sort`
		params = append(params, menu.TYPE_DIR)
		params = append(params, menu.TYPE_MENU)
		params = append(params, userId)
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
func (r *SysMenuImpl) SelectMenuListByRoleId(roleId string, menuCheckStrictly bool) []string {
	querySql := `select m.menu_id as 'str' from sys_menu m 
    left join sys_role_menu rm on m.menu_id = rm.menu_id
    where rm.role_id = ? `
	var params []any
	params = append(params, roleId)
	// 展开
	if menuCheckStrictly {
		querySql += ` and m.menu_id not in 
		(select m.parent_id from sys_menu m 
		inner join sys_role_menu rm on m.menu_id = rm.menu_id 
		and rm.role_id = ?) `
		params = append(params, roleId)
	}

	// 查询结果
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []string{}
	}

	if len(results) > 0 {
		ids := make([]string, 0)
		for _, v := range results {
			ids = append(ids, fmt.Sprintf("%v", v["str"]))
		}
		return ids
	}
	return []string{}
}

// SelectMenuByIds 根据菜单ID查询信息
func (r *SysMenuImpl) SelectMenuByIds(menuIds []string) []model.SysMenu {
	placeholder := datasource.KeyPlaceholderByQuery(len(menuIds))
	querySql := r.selectSql + " where m.menu_id in (" + placeholder + ")"
	parameters := datasource.ConvertIdsSlice(menuIds)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysMenu{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// HasChildByMenuIdAndStatus 存在菜单子节点数量与状态
func (r *SysMenuImpl) HasChildByMenuIdAndStatus(menuId, status string) int64 {
	querySql := "select count(1) as 'total' from sys_menu where parent_id = ?"
	params := []any{menuId}

	// 菜单状态
	if status != "" {
		querySql += " and status = ? and menu_type in (?, ?) "
		params = append(params, status)
		params = append(params, menu.TYPE_DIR)
		params = append(params, menu.TYPE_MENU)
	}

	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return parse.Number(results[0]["total"])
	}
	return 0
}

// InsertMenu 新增菜单信息
func (r *SysMenuImpl) InsertMenu(sysMenu model.SysMenu) string {
	// 参数拼接
	params := make(map[string]any)
	if sysMenu.MenuID != "" {
		params["menu_id"] = sysMenu.MenuID
	}
	if sysMenu.ParentID != "" {
		params["parent_id"] = sysMenu.ParentID
	}
	if sysMenu.MenuName != "" {
		params["menu_name"] = sysMenu.MenuName
	}
	if sysMenu.MenuSort > 0 {
		params["menu_sort"] = sysMenu.MenuSort
	}
	if sysMenu.Path != "" {
		params["path"] = sysMenu.Path
	}
	if sysMenu.Component != "" {
		params["component"] = sysMenu.Component
	}
	if sysMenu.IsFrame != "" {
		params["is_frame"] = sysMenu.IsFrame
	}
	if sysMenu.IsCache != "" {
		params["is_cache"] = sysMenu.IsCache
	}
	if sysMenu.MenuType != "" {
		params["menu_type"] = sysMenu.MenuType
	}
	if sysMenu.Visible != "" {
		params["visible"] = sysMenu.Visible
	}
	if sysMenu.Status != "" {
		params["status"] = sysMenu.Status
	}
	if sysMenu.Perms != "" {
		params["perms"] = sysMenu.Perms
	}
	if sysMenu.Icon != "" {
		params["icon"] = sysMenu.Icon
	} else {
		params["icon"] = "#"
	}
	if sysMenu.Remark != "" {
		params["remark"] = sysMenu.Remark
	}
	if sysMenu.CreateBy != "" {
		params["create_by"] = sysMenu.CreateBy
		params["create_time"] = time.Now().UnixMilli()
	}

	// 根据菜单类型重置参数
	if sysMenu.MenuType == menu.TYPE_BUTTON {
		params["component"] = ""
		params["path"] = ""
		params["icon"] = "#"
		params["is_cache"] = "1"
		params["is_frame"] = "1"
		params["visible"] = "1"
		params["status"] = "1"
	}
	if sysMenu.MenuType == menu.TYPE_DIR {
		params["component"] = ""
		params["perms"] = ""
	}

	// 构建执行语句
	keys, values, placeholder := datasource.KeyValuePlaceholderByInsert(params)
	sql := "insert into sys_menu (" + keys + ")values(" + placeholder + ")"

	db := datasource.DefaultDB()
	// 开启事务
	tx := db.Begin()
	// 执行插入
	err := tx.Exec(sql, values...).Error
	if err != nil {
		logger.Errorf("insert row : %v", err.Error())
		tx.Rollback()
		return ""
	}
	// 获取生成的自增 ID
	var insertedID string
	err = tx.Raw("select last_insert_id()").Row().Scan(&insertedID)
	if err != nil {
		logger.Errorf("insert last id : %v", err.Error())
		tx.Rollback()
		return ""
	}
	// 提交事务
	tx.Commit()
	return insertedID
}

// UpdateMenu 修改菜单信息
func (r *SysMenuImpl) UpdateMenu(sysMenu model.SysMenu) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysMenu.MenuID != "" {
		params["menu_id"] = sysMenu.MenuID
	}
	if sysMenu.ParentID != "" {
		params["parent_id"] = sysMenu.ParentID
	}
	if sysMenu.MenuName != "" {
		params["menu_name"] = sysMenu.MenuName
	}
	if sysMenu.MenuSort > 0 {
		params["menu_sort"] = sysMenu.MenuSort
	}
	if sysMenu.Path != "" {
		params["path"] = sysMenu.Path
	}
	if sysMenu.Component != "" {
		params["component"] = sysMenu.Component
	}
	if sysMenu.IsFrame != "" {
		params["is_frame"] = sysMenu.IsFrame
	}
	if sysMenu.IsCache != "" {
		params["is_cache"] = sysMenu.IsCache
	}
	if sysMenu.MenuType != "" {
		params["menu_type"] = sysMenu.MenuType
	}
	if sysMenu.Visible != "" {
		params["visible"] = sysMenu.Visible
	}
	if sysMenu.Status != "" {
		params["status"] = sysMenu.Status
	}
	if sysMenu.Perms != "" {
		params["perms"] = sysMenu.Perms
	}
	if sysMenu.Icon != "" {
		params["icon"] = sysMenu.Icon
	} else {
		params["icon"] = "#"
	}
	if sysMenu.Remark != "" {
		params["remark"] = sysMenu.Remark
	}
	if sysMenu.UpdateBy != "" {
		params["update_by"] = sysMenu.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 根据菜单类型重置参数
	if sysMenu.MenuType == menu.TYPE_BUTTON {
		params["component"] = ""
		params["path"] = ""
		params["icon"] = "#"
		params["is_cache"] = "1"
		params["is_frame"] = "1"
		params["visible"] = "1"
		params["status"] = "1"
	}
	if sysMenu.MenuType == menu.TYPE_DIR {
		params["component"] = ""
		params["perms"] = ""
	}

	// 构建执行语句
	keys, values := datasource.KeyValueByUpdate(params)
	sql := "update sys_menu set " + keys + " where menu_id = ?"

	// 执行更新
	values = append(values, sysMenu.MenuID)
	rows, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteMenuById 删除菜单管理信息
func (r *SysMenuImpl) DeleteMenuById(menuId string) int64 {
	sql := "delete from sys_menu where menu_id = ?"
	results, err := datasource.ExecDB("", sql, []any{menuId})
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// CheckUniqueMenu 校验菜单是否唯一
func (r *SysMenuImpl) CheckUniqueMenu(sysMenu model.SysMenu) string {
	// 查询条件拼接
	var conditions []string
	var params []any
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
		return ""
	}
	if len(results) > 0 {
		v, ok := results[0]["str"].(string)
		if ok {
			return v
		}
		return ""
	}
	return ""
}
