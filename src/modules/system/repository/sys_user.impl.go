package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	repoUtils "mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strconv"
	"strings"
)

// SysUserImpl 用户表 数据层处理
var SysUserImpl = &sysUserImpl{
	selectSql: `select 
	u.user_id, u.dept_id, u.user_name, u.nick_name, u.user_type, u.email, u.avatar, u.phonenumber, u.password, u.sex, u.status, u.del_flag, u.login_ip, u.login_date, u.create_by, u.create_time, u.remark, 
	d.dept_id, d.parent_id, d.ancestors, d.dept_name, d.order_num, d.leader, d.status as dept_status,
	r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.status as role_status
	from sys_user u
	left join sys_dept d on u.dept_id = d.dept_id
	left join sys_user_role ur on u.user_id = ur.user_id
	left join sys_role r on r.role_id = ur.role_id`,

	sysUserMap: map[string]string{
		"user_id":     "UserID",
		"dept_id":     "DeptID",
		"user_name":   "UserName",
		"nick_name":   "NickName",
		"user_type":   "UserType",
		"email":       "Email",
		"phonenumber": "PhoneNumber",
		"sex":         "Sex",
		"avatar":      "Avatar",
		"password":    "Password",
		"status":      "Status",
		"del_flag":    "DelFlag",
		"login_ip":    "LoginIP",
		"login_date":  "LoginDate",
		"create_by":   "CreateBy",
		"create_time": "CreateTime",
		"update_by":   "UpdateBy",
		"update_time": "UpdateTime",
		"remark":      "Remark",
	},

	sysDeptMap: map[string]string{
		"dept_id":     "DeptID",
		"parent_id":   "ParentID",
		"dept_name":   "DeptName",
		"ancestors":   "Ancestors",
		"order_num":   "OrderNum",
		"leader":      "Leader",
		"dept_status": "Status",
	},

	sysRoleMap: map[string]string{
		"role_id":     "RoleID",
		"role_name":   "RoleName",
		"role_key":    "RoleKey",
		"role_sort":   "RoleSort",
		"data_scope":  "DataScope",
		"role_status": "Status",
	},
}

type sysUserImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 用户信息实体映射
	sysUserMap map[string]string
	// 用户部门实体映射 一对一
	sysDeptMap map[string]string
	// 用户角色实体映射 一对多
	sysRoleMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *sysUserImpl) convertResultRows(rows []map[string]interface{}) []model.SysUser {
	arr := make([]model.SysUser, 0)
	arrKeyIndex := make(map[string]int, 0)

	for i, row := range rows {
		sysUser := model.SysUser{}
		sysDept := model.SysDept{}
		sysRole := model.SysRole{}
		sysUser.Roles = []model.SysRole{}

		for key, value := range row {
			if keyMapper, ok := r.sysUserMap[key]; ok {
				repoUtils.SetFieldValue(&sysUser, keyMapper, value)
			}
			if keyMapper, ok := r.sysDeptMap[key]; ok {
				repoUtils.SetFieldValue(&sysDept, keyMapper, value)
			}
			if keyMapper, ok := r.sysRoleMap[key]; ok {
				repoUtils.SetFieldValue(&sysRole, keyMapper, value)
			}
		}

		sysUser.Dept = sysDept
		if sysRole.RoleKey != "" {
			sysUser.Roles = append(sysUser.Roles, sysRole)
		}

		one := true
		for key, index := range arrKeyIndex {
			if key == sysUser.UserID {
				arrUser := &arr[index]
				arrUser.Roles = append(arrUser.Roles, sysUser.Roles...)
				one = false
				break
			}
		}
		if one {
			arr = append(arr, sysUser)
			arrKeyIndex[sysUser.UserID] = i
		}
	}

	return arr
}

// SelectUserPage 根据条件分页查询用户列表
func (r *sysUserImpl) SelectUserPage(query map[string]string, dataScopeSQL string) map[string]interface{} {
	selectUserSql := `select 
    u.user_id, u.dept_id, u.nick_name, u.user_name, u.email, u.avatar, u.phonenumber, u.sex, u.status, u.del_flag, u.login_ip, u.login_date, u.create_by, u.create_time, u.remark, d.dept_name, d.leader 
    from sys_user u 
	left join sys_dept d on u.dept_id = d.dept_id`
	selectUserTotalSql := `select count(1) as 'total'
    from sys_user u left join sys_dept d on u.dept_id = d.dept_id`

	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if v, ok := query["userId"]; ok {
		conditions = append(conditions, "u.user_id = ?")
		params = append(params, v)
	}
	if v, ok := query["userName"]; ok {
		conditions = append(conditions, "u.user_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok {
		conditions = append(conditions, "u.status = ?")
		params = append(params, v)
	}
	if v, ok := query["phonenumber"]; ok {
		conditions = append(conditions, "u.phonenumber like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["beginTime"]; ok {
		conditions = append(conditions, "u.login_date >= ?")
		beginDate := date.ParseStrToDate(v, date.YYYY_MM_DD)
		params = append(params, beginDate.UnixNano()/1e6)
	}
	if v, ok := query["endTime"]; ok {
		conditions = append(conditions, "u.login_date <= ?")
		endDate := date.ParseStrToDate(v, date.YYYY_MM_DD)
		params = append(params, endDate.UnixNano()/1e6)
	}
	if v, ok := query["deptId"]; ok {
		conditions = append(conditions, "(u.dept_id = ? or u.dept_id in ( select t.dept_id from sys_dept t where find_in_set(?, ancestors) ))")
		params = append(params, v)
		params = append(params, v)
	}

	// 构建查询条件语句
	whereSql := " where u.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询数量 长度为0直接返回
	totalSql := selectUserTotalSql + whereSql + dataScopeSQL
	totalRows, err := datasource.RawDB("", totalSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
	}
	total := parse.Number(totalRows[0]["total"])
	if total <= 0 {
		return map[string]interface{}{
			"total": 0,
			"rows":  []interface{}{},
		}
	}

	// 分页
	pageNum, pageSize := repoUtils.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := selectUserSql + whereSql + dataScopeSQL + pageSql
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
	}

	// 转换实体
	rows := r.convertResultRows(results)
	return map[string]interface{}{
		"total": total,
		"rows":  rows,
	}
}

// SelectAllocatedPage 根据条件分页查询分配用户角色列表
func (r *sysUserImpl) SelectAllocatedPage(query map[string]string, dataScopeSQL string) map[string]interface{} {
	selectUserSql := `select distinct 
    u.user_id, u.dept_id, u.user_name, u.nick_name, u.email, 
    u.phonenumber, u.status, u.create_time, d.dept_name
    from sys_user u
    left join sys_dept d on u.dept_id = d.dept_id
    left join sys_user_role ur on u.user_id = ur.user_id
    left join sys_role r on r.role_id = ur.role_id`
	selectUserTotalSql := `select count(distinct u.user_id) as 'total' from sys_user u
    left join sys_dept d on u.dept_id = d.dept_id
    left join sys_user_role ur on u.user_id = ur.user_id
    left join sys_role r on r.role_id = ur.role_id`

	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if v, ok := query["userName"]; ok {
		conditions = append(conditions, "u.user_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["phonenumber"]; ok {
		conditions = append(conditions, "u.phonenumber like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok {
		conditions = append(conditions, "u.status = ?")
		params = append(params, v)
	}
	// 分配角色用户
	if v, ok := query["allocated"]; ok {
		if parse.Boolean(v) {
			if v, ok := query["deptId"]; ok {
				conditions = append(conditions, "r.role_id = ?")
				params = append(params, v)
			}
		} else {
			if v, ok := query["deptId"]; ok {
				conditions = append(conditions, "(r.role_id != ? or r.role_id IS NULL) and u.user_id not in (select u.user_id from sys_user u inner join sys_user_role ur on u.user_id = ur.user_id and ur.role_id = ?)")
				params = append(params, v)
				params = append(params, v)
			}
		}
	}

	// 构建查询条件语句
	whereSql := " where u.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询数量 长度为0直接返回
	totalSql := selectUserTotalSql + whereSql + dataScopeSQL
	totalRows, err := datasource.RawDB("", totalSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
	}
	total := parse.Number(totalRows[0]["total"])
	if total <= 0 {
		return map[string]interface{}{
			"total": 0,
			"rows":  []interface{}{},
		}
	}

	// 分页
	pageNum, pageSize := repoUtils.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := selectUserSql + whereSql + dataScopeSQL + pageSql
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
	}

	// 转换实体
	rows := r.convertResultRows(results)
	return map[string]interface{}{
		"total": total,
		"rows":  rows,
	}
}

// SelectUserList 根据条件查询用户列表
func (r *sysUserImpl) SelectUserList(sysUser model.SysUser, dataScopeSQL string) []model.SysUser {
	selectUserSql := `select 
    u.user_id, u.dept_id, u.nick_name, u.user_name, u.email, u.avatar, u.phonenumber, u.sex, u.status, u.del_flag, u.login_ip, u.login_date, u.create_by, u.create_time, u.remark, d.dept_name, d.leader 
    from sys_user u
	left join sys_dept d on u.dept_id = d.dept_id`

	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if sysUser.UserID != "" {
		conditions = append(conditions, "u.user_id = ?")
		params = append(params, sysUser.UserID)
	}
	if sysUser.UserName != "" {
		conditions = append(conditions, "u.user_name like concat(?, '%')")
		params = append(params, sysUser.UserName)
	}
	if sysUser.Status != "" {
		conditions = append(conditions, "u.status = ?")
		params = append(params, sysUser.Status)
	}
	if sysUser.PhoneNumber != "" {
		conditions = append(conditions, "u.phonenumber like concat(?, '%')")
		params = append(params, sysUser.PhoneNumber)
	}

	// 构建查询条件语句
	whereSql := " where u.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询数据
	querySql := selectUserSql + whereSql + dataScopeSQL
	rows, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysUser{}
	}
	return r.convertResultRows(rows)
}

// SelectUserByIds 通过用户ID查询用户
func (r *sysUserImpl) SelectUserByIds(userIds []string) []model.SysUser {
	placeholder := repoUtils.KeyPlaceholderByQuery(len(userIds))
	querySql := r.selectSql + " where u.del_flag = '0' and u.user_id in (" + placeholder + ")"
	parameters := repoUtils.ConvertIdsSlice(userIds)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysUser{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// SelectUserByUserName 通过用户登录账号查询用户
func (r *sysUserImpl) SelectUserByUserName(userName string) model.SysUser {
	querySql := r.selectSql + " where u.del_flag = '0' and u.user_name = ?"
	results, err := datasource.RawDB("", querySql, []interface{}{userName})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return model.SysUser{}
	}
	// 转换实体
	rows := r.convertResultRows(results)
	if len(rows) > 0 {
		return rows[0]
	}
	return model.SysUser{}
}

// InsertUser 新增用户信息
func (r *sysUserImpl) InsertUser(sysUser model.SysUser) string {
	return ""
}

// UpdateUser 修改用户信息
func (r *sysUserImpl) UpdateUser(sysUser model.SysUser) int {
	return 0
}

// DeleteUserByIds 批量删除用户信息
func (r *sysUserImpl) DeleteUserByIds(userIds []string) int64 {
	placeholder := repoUtils.KeyPlaceholderByQuery(len(userIds))
	sql := "update sys_user set del_flag = '1' where user_id in (" + placeholder + ")"
	parameters := repoUtils.ConvertIdsSlice(userIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("update err => %v", err)
		return 0
	}
	return results
}

// CheckUniqueUser 校验用户信息是否唯一
func (r *sysUserImpl) CheckUniqueUser(sysUser model.SysUser) string {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if sysUser.UserName != "" {
		conditions = append(conditions, "user_name = ?")
		params = append(params, sysUser.UserName)
	}
	if sysUser.PhoneNumber != "" {
		conditions = append(conditions, "phonenumber = ?")
		params = append(params, sysUser.PhoneNumber)
	}
	if sysUser.Email != "" {
		conditions = append(conditions, "email = ?")
		params = append(params, sysUser.Email)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	} else {
		return ""
	}

	// 查询数据
	querySql := "select user_id as 'str' from sys_user" + whereSql + " limit 1"
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
	}
	if len(results) > 0 {
		return strconv.FormatInt(results[0]["str"].(int64), 10)
	}
	return ""
}
