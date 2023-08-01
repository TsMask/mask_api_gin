package repository

import (
	"fmt"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// 实例化数据层 SysRoleImpl 结构体
var NewSysRoleImpl = &SysRoleImpl{
	selectSql: `select distinct 
	r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.menu_check_strictly, 
	r.dept_check_strictly, r.status, r.del_flag, r.create_time, r.remark 
	from sys_role r
	left join sys_user_role ur on ur.role_id = r.role_id
	left join sys_user u on u.user_id = ur.user_id
	left join sys_dept d on u.dept_id = d.dept_id`,

	resultMap: map[string]string{
		"role_id":             "RoleID",
		"role_name":           "RoleName",
		"role_key":            "RoleKey",
		"role_sort":           "RoleSort",
		"data_scope":          "DataScope",
		"menu_check_strictly": "MenuCheckStrictly",
		"dept_check_strictly": "DeptCheckStrictly",
		"status":              "Status",
		"del_flag":            "DelFlag",
		"create_by":           "CreateBy",
		"create_time":         "CreateTime",
		"update_by":           "UpdateBy",
		"update_time":         "UpdateTime",
		"remark":              "Remark",
	},
}

// SysRoleImpl 角色表 数据层处理
type SysRoleImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysRoleImpl) convertResultRows(rows []map[string]any) []model.SysRole {
	arr := make([]model.SysRole, 0)
	for _, row := range rows {
		sysRole := model.SysRole{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysRole, keyMapper, value)
			}
		}
		arr = append(arr, sysRole)
	}
	return arr
}

// SelectRolePage 根据条件分页查询角色数据
func (r *SysRoleImpl) SelectRolePage(query map[string]any, dataScopeSQL string) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["roleId"]; ok && v != "" {
		conditions = append(conditions, "r.role_id = ?")
		params = append(params, v)
	}
	if v, ok := query["roleName"]; ok && v != "" {
		conditions = append(conditions, "r.role_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["roleKey"]; ok && v != "" {
		conditions = append(conditions, "r.role_key like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok && v != "" {
		conditions = append(conditions, "r.status = ?")
		params = append(params, v)
	}
	if v, ok := query["beginTime"]; ok && v != "" {
		conditions = append(conditions, "r.create_time >= ?")
		beginDate := date.ParseStrToDate(v.(string), date.YYYY_MM_DD)
		params = append(params, beginDate.UnixMilli())
	}
	if v, ok := query["endTime"]; ok && v != "" {
		conditions = append(conditions, "r.create_time <= ?")
		endDate := date.ParseStrToDate(v.(string), date.YYYY_MM_DD)
		params = append(params, endDate.UnixMilli())
	}
	if v, ok := query["deptId"]; ok && v != "" {
		conditions = append(conditions, `(u.dept_id = ? or u.dept_id in ( 
			select t.dept_id from sys_dept t where find_in_set(?, ancestors)
		))`)
		params = append(params, v)
		params = append(params, v)
	}

	// 构建查询条件语句
	whereSql := " where r.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询数量 长度为0直接返回
	totalSql := `select count(distinct r.role_id) as 'total' from sys_role r
    left join sys_user_role ur on ur.role_id = r.role_id
    left join sys_user u on u.user_id = ur.user_id
    left join sys_dept d on u.dept_id = d.dept_id`
	totalRows, err := datasource.RawDB("", totalSql+whereSql+dataScopeSQL, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
	}
	total := parse.Number(totalRows[0]["total"])
	if total == 0 {
		return map[string]any{
			"total": total,
			"rows":  []model.SysRole{},
		}
	}

	// 分页
	pageNum, pageSize := repo.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " order by r.role_sort asc limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := r.selectSql + whereSql + dataScopeSQL + pageSql
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
	}

	// 转换实体
	rows := r.convertResultRows(results)
	return map[string]any{
		"total": total,
		"rows":  rows,
	}
}

// SelectRoleList 根据条件查询角色数据
func (r *SysRoleImpl) SelectRoleList(sysRole model.SysRole, dataScopeSQL string) []model.SysRole {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysRole.RoleID != "" {
		conditions = append(conditions, "r.role_id = ?")
		params = append(params, sysRole.RoleID)
	}
	if sysRole.RoleKey != "" {
		conditions = append(conditions, "r.role_key like concat(?, '%')")
		params = append(params, sysRole.RoleKey)
	}
	if sysRole.RoleName != "" {
		conditions = append(conditions, "r.role_name like concat(?, '%')")
		params = append(params, sysRole.RoleName)
	}
	if sysRole.Status != "" {
		conditions = append(conditions, "r.status = ?")
		params = append(params, sysRole.Status)
	}

	// 构建查询条件语句
	whereSql := " where r.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询数据
	orderSql := " order by r.role_sort"
	querySql := r.selectSql + whereSql + dataScopeSQL + orderSql
	rows, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysRole{}
	}
	return r.convertResultRows(rows)
}

// SelectRoleListByUserId 根据用户ID获取角色选择框列表
func (r *SysRoleImpl) SelectRoleListByUserId(userId string) []model.SysRole {
	querySql := r.selectSql + " where r.del_flag = '0' and ur.user_id = ?"
	results, err := datasource.RawDB("", querySql, []any{userId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysRole{}
	}
	return r.convertResultRows(results)
}

// SelectRoleByIds 通过角色ID查询角色
func (r *SysRoleImpl) SelectRoleByIds(roleIds []string) []model.SysRole {
	placeholder := repo.KeyPlaceholderByQuery(len(roleIds))
	querySql := r.selectSql + " where r.role_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(roleIds)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysRole{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// UpdateRole 修改角色信息
func (r *SysRoleImpl) UpdateRole(sysRole model.SysRole) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysRole.RoleName != "" {
		params["role_name"] = sysRole.RoleName
	}
	if sysRole.RoleKey != "" {
		params["role_key"] = sysRole.RoleKey
	}
	if sysRole.RoleSort > 0 {
		params["role_sort"] = sysRole.RoleSort
	}
	if sysRole.DataScope != "" {
		params["data_scope"] = sysRole.DataScope
	}
	if sysRole.MenuCheckStrictly != "" {
		params["menu_check_strictly"] = sysRole.MenuCheckStrictly
	}
	if sysRole.DeptCheckStrictly != "" {
		params["dept_check_strictly"] = sysRole.DeptCheckStrictly
	}
	if sysRole.Status != "" {
		params["status"] = sysRole.Status
	}
	if sysRole.Remark != "" {
		params["remark"] = sysRole.Remark
	}
	if sysRole.UpdateBy != "" {
		params["update_by"] = sysRole.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := repo.KeyValueByUpdate(params)
	sql := "update sys_role set " + strings.Join(keys, ",") + " where role_id = ?"

	// 执行更新
	values = append(values, sysRole.RoleID)
	rows, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// InsertRole 新增角色信息
func (r *SysRoleImpl) InsertRole(sysRole model.SysRole) string {
	// 参数拼接
	params := make(map[string]any)
	if sysRole.RoleID != "" {
		params["role_id"] = sysRole.RoleID
	}
	if sysRole.RoleName != "" {
		params["role_name"] = sysRole.RoleName
	}
	if sysRole.RoleKey != "" {
		params["role_key"] = sysRole.RoleKey
	}
	if sysRole.RoleSort > 0 {
		params["role_sort"] = sysRole.RoleSort
	}
	if sysRole.DataScope != "" {
		params["data_scope"] = sysRole.DataScope
	}
	if sysRole.MenuCheckStrictly != "" {
		params["menu_check_strictly"] = sysRole.MenuCheckStrictly
	}
	if sysRole.DeptCheckStrictly != "" {
		params["dept_check_strictly"] = sysRole.DeptCheckStrictly
	}
	if sysRole.Status != "" {
		params["status"] = sysRole.Status
	}
	if sysRole.Remark != "" {
		params["remark"] = sysRole.Remark
	}
	if sysRole.CreateBy != "" {
		params["create_by"] = sysRole.CreateBy
		params["create_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_role (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// DeleteRoleByIds 批量删除角色信息
func (r *SysRoleImpl) DeleteRoleByIds(roleIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(roleIds))
	sql := "update sys_role set del_flag = '1' where role_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(roleIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// CheckUniqueRole 校验角色是否唯一
func (r *SysRoleImpl) CheckUniqueRole(sysRole model.SysRole) string {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysRole.RoleName != "" {
		conditions = append(conditions, "r.role_name = ?")
		params = append(params, sysRole.RoleName)
	}
	if sysRole.RoleKey != "" {
		conditions = append(conditions, "r.role_key = ?")
		params = append(params, sysRole.RoleKey)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	} else {
		return ""
	}

	// 查询数据
	querySql := "select role_id as 'str' from sys_role r " + whereSql + " and r.del_flag = '0' limit 1"
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
		return ""
	}
	if len(results) > 0 {
		return fmt.Sprintf("%v", results[0]["str"])
	}
	return ""
}
