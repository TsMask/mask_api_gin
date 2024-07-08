package repository

import (
	"fmt"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/crypto"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// NewSysUser 实例化数据层
var NewSysUser = &SysUserRepository{
	selectSql: `select 
	u.user_id, u.dept_id, u.user_name, u.nick_name, u.user_type, u.email, u.avatar, u.phone, u.password, u.sex, u.status, u.del_flag, u.login_ip, u.login_date, u.create_by, u.create_time, u.remark, 
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
		"phone":       "Phone",
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

// SysUserRepository 用户表 数据层处理
type SysUserRepository struct {
	selectSql  string            // 查询视图对象SQL
	sysUserMap map[string]string // 用户信息实体映射
	sysDeptMap map[string]string // 用户部门实体映射 一对一
	sysRoleMap map[string]string // 用户角色实体映射 一对多
}

// convertResultRows 将结果记录转实体结果组
func (r *SysUserRepository) convertResultRows(rows []map[string]any) []model.SysUser {
	arr := make([]model.SysUser, 0)

	for _, row := range rows {
		sysUser := model.SysUser{}
		sysDept := model.SysDept{}
		sysRole := model.SysRole{}
		sysUser.Roles = []model.SysRole{}

		for key, value := range row {
			if keyMapper, ok := r.sysUserMap[key]; ok {
				db.SetFieldValue(&sysUser, keyMapper, value)
			}
			if keyMapper, ok := r.sysDeptMap[key]; ok {
				db.SetFieldValue(&sysDept, keyMapper, value)
			}
			if keyMapper, ok := r.sysRoleMap[key]; ok {
				db.SetFieldValue(&sysRole, keyMapper, value)
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

// SelectByPage 分页查询集合
func (r *SysUserRepository) SelectByPage(queryMap map[string]any, dataScopeSQL string) map[string]any {
	selectUserSql := `select 
    u.user_id, u.dept_id, u.nick_name, u.user_name, u.email, u.avatar, u.phone, u.sex, u.status, u.del_flag, u.login_ip, u.login_date, u.create_by, u.create_time, u.remark, d.dept_name, d.leader 
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
	if v, ok := queryMap["phone"]; ok && v != "" {
		conditions = append(conditions, "u.phone like concat(?, '%')")
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
		"total": int64(0),
		"rows":  []model.SysUser{},
	}

	// 查询数量 长度为0直接返回
	totalSql := selectUserTotalSql + whereSql + dataScopeSQL
	totalRows, err := db.RawDB("", totalSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
		return result
	}
	if total := parse.Number(totalRows[0]["total"]); total > 0 {
		result["total"] = total
	} else {
		return result
	}

	// 分页
	pageNum, pageSize := db.PageNumSize(queryMap["pageNum"], queryMap["pageSize"])
	pageSql := " limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := selectUserSql + whereSql + dataScopeSQL + pageSql
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return result
	}

	// 转换实体
	result["rows"] = r.convertResultRows(results)
	return result
}

// Select 查询集合
func (r *SysUserRepository) Select(sysUser model.SysUser, dataScopeSQL string) []model.SysUser {
	selectUserSql := `select 
    u.user_id, u.dept_id, u.nick_name, u.user_name, u.email, u.avatar, u.phone, u.sex, u.status, 
    u.del_flag, u.login_ip, u.login_date, u.create_by, u.create_time, u.remark, 
    d.dept_name, d.leader 
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
	if sysUser.Phone != "" {
		conditions = append(conditions, "u.phone like concat(?, '%')")
		params = append(params, sysUser.Phone)
	}

	// 构建查询条件语句
	whereSql := " where u.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询数据
	querySql := selectUserSql + whereSql + dataScopeSQL
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysUser{}
	}
	return r.convertResultRows(rows)
}

// SelectByIds 通过ID查询信息
func (r *SysUserRepository) SelectByIds(userIds []string) []model.SysUser {
	placeholder := db.KeyPlaceholderByQuery(len(userIds))
	querySql := r.selectSql + " where u.del_flag = '0' and u.user_id in (" + placeholder + ")"
	parameters := db.ConvertIdsSlice(userIds)
	results, err := db.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysUser{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// Insert 新增信息
func (r *SysUserRepository) Insert(sysUser model.SysUser) string {
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
	if sysUser.Phone != "" {
		params["phone"] = sysUser.Phone
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
	params["remark"] = sysUser.Remark
	if sysUser.CreateBy != "" {
		params["create_by"] = sysUser.CreateBy
		params["create_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_user (%s)values(%s)", keys, placeholder)

	tx := db.DB("").Begin() // 开启事务
	// 执行插入
	if err := tx.Exec(sql, values...).Error; err != nil {
		logger.Errorf("insert row : %v", err.Error())
		tx.Rollback()
		return ""
	}
	// 获取生成的自增ID
	var insertedID string
	if err := tx.Raw("select last_insert_id()").Row().Scan(&insertedID); err != nil {
		logger.Errorf("insert last id : %v", err.Error())
		tx.Rollback()
		return ""
	}
	tx.Commit() // 提交事务
	return insertedID
}

// Update 修改信息
func (r *SysUserRepository) Update(sysUser model.SysUser) int64 {
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
	params["email"] = sysUser.Email
	params["phone"] = sysUser.Phone
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
	keys, values := db.KeyValueByUpdate(params)
	sql := fmt.Sprintf("update sys_user set %s where user_id = ?", keys)

	// 执行更新
	values = append(values, sysUser.UserID)
	rows, err := db.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteByIds 批量删除信息
func (r *SysUserRepository) DeleteByIds(userIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(userIds))
	sql := "update sys_user user_name = concat(user_name, '_del'), set del_flag = '1' where user_id in (" + placeholder + ")"
	parameters := db.ConvertIdsSlice(userIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("update err => %v", err)
		return 0
	}
	return results
}

// CheckUnique 检查信息是否唯一
func (r *SysUserRepository) CheckUnique(sysUser model.SysUser) string {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysUser.UserName != "" {
		conditions = append(conditions, "user_name = ?")
		params = append(params, sysUser.UserName)
	}
	if sysUser.Phone != "" {
		conditions = append(conditions, "phone = ?")
		params = append(params, sysUser.Phone)
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
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
	}
	if len(results) > 0 {
		return fmt.Sprint(results[0]["str"])
	}
	return ""
}

// SelectAllocatedByPage 分页查询集合By分配用户角色
func (r *SysUserRepository) SelectAllocatedByPage(query map[string]any, dataScopeSQL string) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["userName"]; ok && v != "" {
		conditions = append(conditions, "u.user_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["phone"]; ok && v != "" {
		conditions = append(conditions, "u.phone like concat(?, '%')")
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
		"total": int64(0),
		"rows":  []model.SysUser{},
	}

	// 查询数量 长度为0直接返回
	totalSql := `select count(distinct u.user_id) as 'total' from sys_user u
    left join sys_dept d on u.dept_id = d.dept_id
    left join sys_user_role ur on u.user_id = ur.user_id
    left join sys_role r on r.role_id = ur.role_id`
	totalRows, err := db.RawDB("", totalSql+whereSql+dataScopeSQL, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
		return result
	}
	if total := parse.Number(totalRows[0]["total"]); total > 0 {
		result["total"] = total
	} else {
		return result
	}

	// 分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := `select distinct 
    u.user_id, u.dept_id, u.user_name, u.nick_name, u.email, 
    u.phone, u.status, u.create_time, d.dept_name
    from sys_user u
    left join sys_dept d on u.dept_id = d.dept_id
    left join sys_user_role ur on u.user_id = ur.user_id
    left join sys_role r on r.role_id = ur.role_id`
	querySql = querySql + whereSql + dataScopeSQL + pageSql
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
	}

	// 转换实体
	result["rows"] = r.convertResultRows(results)
	return result
}

// SelectByUserName 通过登录账号查询信息
func (r *SysUserRepository) SelectByUserName(userName string) model.SysUser {
	querySql := r.selectSql + " where u.del_flag = '0' and u.user_name = ?"
	results, err := db.RawDB("", querySql, []any{userName})
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
