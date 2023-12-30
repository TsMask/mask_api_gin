package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/crypto"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// 实例化数据层 SysUserImpl 结构体
var NewSysUserImpl = &SysUserImpl{
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

// SysUserImpl 用户表 数据层处理
type SysUserImpl struct {
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
func (r *SysUserImpl) convertResultRows(rows []map[string]any) []model.SysUser {
	arr := make([]model.SysUser, 0)

	for _, row := range rows {
		sysUser := model.SysUser{}
		sysDept := model.SysDept{}
		sysRole := model.SysRole{}
		sysUser.Roles = []model.SysRole{}

		for key, value := range row {
			if keyMapper, ok := r.sysUserMap[key]; ok {
				repo.SetFieldValue(&sysUser, keyMapper, value)
			}
			if keyMapper, ok := r.sysDeptMap[key]; ok {
				repo.SetFieldValue(&sysDept, keyMapper, value)
			}
			if keyMapper, ok := r.sysRoleMap[key]; ok {
				repo.SetFieldValue(&sysRole, keyMapper, value)
			}
		}

		sysUser.Dept = sysDept
		if sysRole.RoleKey != "" {
			sysUser.Roles = append(sysUser.Roles, sysRole)
		}

		one := true
		for i, a := range arr {
			if a.UserID == sysUser.UserID {
				arrUser := &arr[i]
				arrUser.Roles = append(arrUser.Roles, sysUser.Roles...)
				one = false
				break
			}
		}
		if one {
			arr = append(arr, sysUser)
		}
	}

	return arr
}

// SelectUserPage 根据条件分页查询用户列表
func (r *SysUserImpl) SelectUserPage(queryMap map[string]any, dataScopeSQL string) map[string]any {
	selectUserSql := `select 
    u.user_id, u.dept_id, u.nick_name, u.user_name, u.email, u.avatar, u.phonenumber, u.sex, u.status, u.del_flag, u.login_ip, u.login_date, u.create_by, u.create_time, u.remark, d.dept_name, d.leader 
    from sys_user u 
	left join sys_dept d on u.dept_id = d.dept_id`
	selectUserTotalSql := `select count(distinct u.user_id) as 'total'
    from sys_user u left join sys_dept d on u.dept_id = d.dept_id`

	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := queryMap["userId"]; ok && v != "" {
		conditions = append(conditions, "u.user_id = ?")
		params = append(params, v)
	}
	if v, ok := queryMap["userName"]; ok && v != "" {
		conditions = append(conditions, "u.user_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := queryMap["status"]; ok && v != "" {
		conditions = append(conditions, "u.status = ?")
		params = append(params, v)
	}
	if v, ok := queryMap["phonenumber"]; ok && v != "" {
		conditions = append(conditions, "u.phonenumber like concat(?, '%')")
		params = append(params, v)
	}
	beginTime, ok := queryMap["beginTime"]
	if !ok {
		beginTime, ok = queryMap["params[beginTime]"]
	}
	if ok && beginTime != "" {
		conditions = append(conditions, "u.login_date >= ?")
		beginDate := date.ParseStrToDate(beginTime.(string), date.YYYY_MM_DD)
		params = append(params, beginDate.UnixMilli())
	}
	endTime, ok := queryMap["endTime"]
	if !ok {
		endTime, ok = queryMap["params[endTime]"]
	}
	if ok && endTime != "" {
		conditions = append(conditions, "u.login_date <= ?")
		endDate := date.ParseStrToDate(endTime.(string), date.YYYY_MM_DD)
		params = append(params, endDate.UnixMilli())
	}
	if v, ok := queryMap["deptId"]; ok && v != "" {
		conditions = append(conditions, "(u.dept_id = ? or u.dept_id in ( select t.dept_id from sys_dept t where find_in_set(?, ancestors) ))")
		params = append(params, v)
		params = append(params, v)
	}

	// 构建查询条件语句
	whereSql := " where u.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询结果
	result := map[string]any{
		"total": 0,
		"rows":  []model.SysUser{},
	}

	// 查询数量 长度为0直接返回
	totalSql := selectUserTotalSql + whereSql + dataScopeSQL
	totalRows, err := datasource.RawDB("", totalSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
		return result
	}
	total := parse.Number(totalRows[0]["total"])
	if total == 0 {
		return result
	} else {
		result["total"] = total
	}

	// 分页
	pageNum, pageSize := repo.PageNumSize(queryMap["pageNum"], queryMap["pageSize"])
	pageSql := " limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := selectUserSql + whereSql + dataScopeSQL + pageSql
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return result
	}

	// 转换实体
	result["rows"] = r.convertResultRows(results)
	return result
}

// SelectAllocatedPage 根据条件分页查询分配用户角色列表
func (r *SysUserImpl) SelectAllocatedPage(query map[string]any, dataScopeSQL string) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["userName"]; ok && v != "" {
		conditions = append(conditions, "u.user_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["phonenumber"]; ok && v != "" {
		conditions = append(conditions, "u.phonenumber like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok && v != "" {
		conditions = append(conditions, "u.status = ?")
		params = append(params, v)
	}
	// 分配角色用户
	if allocated, ok := query["allocated"]; ok && allocated != "" {
		if roleId, ok := query["roleId"]; ok && roleId != "" {
			if parse.Boolean(allocated) {
				conditions = append(conditions, "r.role_id = ?")
				params = append(params, roleId)
			} else {
				conditions = append(conditions, `(r.role_id != ? or r.role_id IS NULL) 
				and u.user_id not in (
					select u.user_id from sys_user u 
					inner join sys_user_role ur on u.user_id = ur.user_id 
					and ur.role_id = ?
				)`)
				params = append(params, roleId)
				params = append(params, roleId)
			}

		}
	}

	// 构建查询条件语句
	whereSql := " where u.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询结果
	result := map[string]any{
		"total": 0,
		"rows":  []model.SysUser{},
	}

	// 查询数量 长度为0直接返回
	totalSql := `select count(distinct u.user_id) as 'total' from sys_user u
    left join sys_dept d on u.dept_id = d.dept_id
    left join sys_user_role ur on u.user_id = ur.user_id
    left join sys_role r on r.role_id = ur.role_id`
	totalRows, err := datasource.RawDB("", totalSql+whereSql+dataScopeSQL, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
		return result
	}
	total := parse.Number(totalRows[0]["total"])
	if total == 0 {
		return result
	} else {
		result["total"] = total
	}

	// 分页
	pageNum, pageSize := repo.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := `select distinct 
    u.user_id, u.dept_id, u.user_name, u.nick_name, u.email, 
    u.phonenumber, u.status, u.create_time, d.dept_name
    from sys_user u
    left join sys_dept d on u.dept_id = d.dept_id
    left join sys_user_role ur on u.user_id = ur.user_id
    left join sys_role r on r.role_id = ur.role_id`
	querySql = querySql + whereSql + dataScopeSQL + pageSql
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
	}

	// 转换实体
	result["rows"] = r.convertResultRows(results)
	return result
}

// SelectUserList 根据条件查询用户列表
func (r *SysUserImpl) SelectUserList(sysUser model.SysUser, dataScopeSQL string) []model.SysUser {
	selectUserSql := `select 
    u.user_id, u.dept_id, u.nick_name, u.user_name, u.email, u.avatar, u.phonenumber, u.sex, u.status, u.del_flag, u.login_ip, u.login_date, u.create_by, u.create_time, u.remark, d.dept_name, d.leader 
    from sys_user u
	left join sys_dept d on u.dept_id = d.dept_id`

	// 查询条件拼接
	var conditions []string
	var params []any
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
func (r *SysUserImpl) SelectUserByIds(userIds []string) []model.SysUser {
	placeholder := repo.KeyPlaceholderByQuery(len(userIds))
	querySql := r.selectSql + " where u.del_flag = '0' and u.user_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(userIds)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysUser{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// SelectUserByUserName 通过用户登录账号查询用户
func (r *SysUserImpl) SelectUserByUserName(userName string) model.SysUser {
	querySql := r.selectSql + " where u.del_flag = '0' and u.user_name = ?"
	results, err := datasource.RawDB("", querySql, []any{userName})
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
func (r *SysUserImpl) InsertUser(sysUser model.SysUser) string {
	// 参数拼接
	params := make(map[string]any)
	if sysUser.UserID != "" {
		params["user_id"] = sysUser.UserID
	}
	if sysUser.DeptID != "" {
		params["dept_id"] = sysUser.DeptID
	}
	if sysUser.UserName != "" {
		params["user_name"] = sysUser.UserName
	}
	if sysUser.NickName != "" {
		params["nick_name"] = sysUser.NickName
	}
	if sysUser.UserType != "" {
		params["user_type"] = sysUser.UserType
	}
	if sysUser.Avatar != "" {
		params["avatar"] = sysUser.Avatar
	}
	if sysUser.Email != "" {
		params["email"] = sysUser.Email
	}
	if sysUser.PhoneNumber != "" {
		params["phonenumber"] = sysUser.PhoneNumber
	}
	if sysUser.Sex != "" {
		params["sex"] = sysUser.Sex
	}
	if sysUser.Password != "" {
		password := crypto.BcryptHash(sysUser.Password)
		params["password"] = password
	}
	if sysUser.Status != "" {
		params["status"] = sysUser.Status
	}
	if sysUser.Remark != "" {
		params["remark"] = sysUser.Remark
	}
	if sysUser.CreateBy != "" {
		params["create_by"] = sysUser.CreateBy
		params["create_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_user (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// UpdateUser 修改用户信息
func (r *SysUserImpl) UpdateUser(sysUser model.SysUser) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysUser.DeptID != "" {
		params["dept_id"] = sysUser.DeptID
	}
	if sysUser.UserName != "" {
		params["user_name"] = sysUser.UserName
	}
	if sysUser.NickName != "" {
		params["nick_name"] = sysUser.NickName
	}
	if sysUser.UserType != "" {
		params["user_type"] = sysUser.UserType
	}
	if sysUser.Avatar != "" {
		params["avatar"] = sysUser.Avatar
	}
	if sysUser.Email != "" {
		if sysUser.Email == "nil" {
			params["email"] = ""
		} else {
			params["email"] = sysUser.Email
		}
	}
	if sysUser.PhoneNumber != "" {
		if sysUser.PhoneNumber == "nil" {
			params["phonenumber"] = ""
		} else {
			params["phonenumber"] = sysUser.PhoneNumber
		}
	}
	if sysUser.Sex != "" {
		params["sex"] = sysUser.Sex
	}
	if sysUser.Password != "" {
		password := crypto.BcryptHash(sysUser.Password)
		params["password"] = password
	}
	if sysUser.Status != "" {
		params["status"] = sysUser.Status
	}
	if sysUser.Remark != "" {
		params["remark"] = sysUser.Remark
	}
	if sysUser.UpdateBy != "" {
		params["update_by"] = sysUser.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}
	if sysUser.LoginIP != "" {
		params["login_ip"] = sysUser.LoginIP
	}
	if sysUser.LoginDate > 0 {
		params["login_date"] = sysUser.LoginDate
	}

	// 构建执行语句
	keys, values := repo.KeyValueByUpdate(params)
	sql := "update sys_user set " + strings.Join(keys, ",") + " where user_id = ?"

	// 执行更新
	values = append(values, sysUser.UserID)
	rows, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteUserByIds 批量删除用户信息
func (r *SysUserImpl) DeleteUserByIds(userIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(userIds))
	sql := "update sys_user set del_flag = '1' where user_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(userIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("update err => %v", err)
		return 0
	}
	return results
}

// CheckUniqueUser 校验用户信息是否唯一
func (r *SysUserImpl) CheckUniqueUser(sysUser model.SysUser) string {
	// 查询条件拼接
	var conditions []string
	var params []any
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
	querySql := "select user_id as 'str' from sys_user " + whereSql + " limit 1"
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
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
