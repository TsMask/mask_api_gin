package repository

import (
	"fmt"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// NewSysRole 实例化数据层
var NewSysRole = &SysRole{
	selectSql: `select distinct 
	r.role_id, r.role_name, r.role_key, r.role_sort, 
	r.data_scope, r.menu_check_strictly, r.dept_check_strictly, 
	r.status, r.del_flag, r.create_time, r.remark 
	from sys_role r`,

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

// SysRole 角色表 数据层处理
type SysRole struct {
	selectSql string            // 查询视图对象SQL
	resultMap map[string]string // 结果字段与实体映射
}

// SelectByPage 分页查询集合
func (r SysRole) SelectByPage(query map[string]any, dataScopeSQL string) ([]model.SysRole, int64) {
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
	beginTime, ok := query["beginTime"]
	if !ok {
		beginTime, ok = query["params[beginTime]"]
	}
	if ok && beginTime != "" {
		conditions = append(conditions, "r.create_time >= ?")
		beginDate := date.ParseStrToDate(beginTime.(string), date.YYYY_MM_DD)
		params = append(params, beginDate.UnixMilli())
	}
	endTime, ok := query["endTime"]
	if !ok {
		endTime, ok = query["params[endTime]"]
	}
	if ok && endTime != "" {
		conditions = append(conditions, "r.create_time <= ?")
		endDate := date.ParseStrToDate(endTime.(string), date.YYYY_MM_DD)
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

	// 查询结果
	total := int64(0)
	arr := []model.SysRole{}

	// 查询数量 长度为0直接返回
	totalSql := `select count(distinct r.role_id) as 'total' from sys_role r
    left join sys_user_role ur on ur.role_id = r.role_id
    left join sys_user u on u.user_id = ur.user_id
    left join sys_dept d on u.dept_id = d.dept_id`
	totalRows, err := db.RawDB("", totalSql+whereSql+dataScopeSQL, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
		return arr, total
	}
	total = parse.Number(totalRows[0]["total"])
	if total <= 0 {
		return arr, total
	}

	// 分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " order by r.role_sort asc limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := r.selectSql + whereSql + dataScopeSQL + pageSql
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return arr, total
	}

	// 转换实体
	if err := db.Unmarshal(rows, &arr); err != nil {
		logger.Errorf("unmarshal err => %v", err)
	}
	return arr, total
}

// Select 查询集合
func (r SysRole) Select(sysRole model.SysRole, dataScopeSQL string) []model.SysRole {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysRole.RoleId != "" {
		conditions = append(conditions, "r.role_id = ?")
		params = append(params, sysRole.RoleId)
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
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysRole{}
	}
	// 转换实体
	arr := []model.SysRole{}
	if err := db.Unmarshal(rows, &arr); err != nil {
		logger.Errorf("unmarshal err => %v", err)
	}
	return arr
}

// SelectByIds 通过ID查询信息
func (r SysRole) SelectByIds(roleIds []string) []model.SysRole {
	placeholder := db.KeyPlaceholderByQuery(len(roleIds))
	querySql := r.selectSql + " where r.role_id in (" + placeholder + ")"
	parameters := db.ConvertIdsSlice(roleIds)
	rows, err := db.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysRole{}
	}
	// 转换实体
	arr := []model.SysRole{}
	if err := db.Unmarshal(rows, &arr); err != nil {
		logger.Errorf("unmarshal err => %v", err)
	}
	return arr
}

// Update 修改信息
func (r SysRole) Update(sysRole model.SysRole) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysRole.RoleName != "" {
		params["role_name"] = sysRole.RoleName
	}
	if sysRole.RoleKey != "" {
		params["role_key"] = sysRole.RoleKey
	}
	if sysRole.RoleSort >= 0 {
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
	params["remark"] = sysRole.Remark
	if sysRole.UpdateBy != "" {
		params["update_by"] = sysRole.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := db.KeyValueByUpdate(params)
	sql := fmt.Sprintf("update sys_role set %s where role_id = ?", keys)

	// 执行更新
	values = append(values, sysRole.RoleId)
	rows, err := db.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// Insert 新增信息
func (r SysRole) Insert(sysRole model.SysRole) string {
	// 参数拼接
	params := make(map[string]any)
	if sysRole.RoleId != "" {
		params["role_id"] = sysRole.RoleId
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
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_role (%s)values(%s)", keys, placeholder)

	tx := db.DB("").Begin() // 开启事务
	// 执行插入
	if err := tx.Exec(sql, values...).Error; err != nil {
		logger.Errorf("insert row : %v", err.Error())
		tx.Rollback()
		return ""
	}
	// 获取生成的自增 ID
	var insertedID string
	if err := tx.Raw("select last_insert_id()").Row().Scan(&insertedID); err != nil {
		logger.Errorf("insert last id : %v", err.Error())
		tx.Rollback()
		return ""
	}
	tx.Commit() // 提交事务
	return insertedID
}

// DeleteByIds 批量删除信息
func (r SysRole) DeleteByIds(roleIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(roleIds))
	sql := fmt.Sprintf("update sys_role set del_flag = '1' where role_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(roleIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// SelectByUserId 根据用户ID获取角色信息
func (r SysRole) SelectByUserId(userId string) []model.SysRole {
	querySql := `select 
	r.role_id, r.role_name, r.role_key, r.role_sort, 
	r.data_scope, r.menu_check_strictly, r.dept_check_strictly, 
	r.status, r.del_flag, r.create_time, r.remark 
	from sys_user_role ur 
	left join sys_user u on u.user_id = ur.user_id
	left join sys_role r on r.role_id = ur.role_id
	where u.del_flag = '0' and ur.user_id = ?`
	rows, err := db.RawDB("", querySql, []any{userId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysRole{}
	}
	// 转换实体
	arr := []model.SysRole{}
	if err := db.Unmarshal(rows, &arr); err != nil {
		logger.Errorf("unmarshal err => %v", err)
	}
	return arr
}

// CheckUnique 检查信息是否唯一
func (r SysRole) CheckUnique(sysRole model.SysRole) string {
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
		return "-"
	}

	// 查询数据
	querySql := "select role_id as 'str' from sys_role r " + whereSql + " and r.del_flag = '0' limit 1"
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
		return "-"
	}
	if len(results) > 0 {
		return fmt.Sprint(results[0]["str"])
	}
	return ""
}
