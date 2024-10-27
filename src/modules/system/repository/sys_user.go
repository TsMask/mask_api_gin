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
var NewSysUser = &SysUser{
	sql: `select 
	u.user_id, u.dept_id, u.user_name, u.nick_name, u.user_type, u.email, u.avatar, 
	u.phone, u.password, u.sex, u.status, u.del_flag, u.login_ip, u.login_date, 
	u.create_by, u.create_time, u.remark
	from sys_user u`,
}

// SysUser 用户表 数据层处理
type SysUser struct {
	sql string // 查询视图对象SQL
}

// SelectByPage 分页查询集合
func (r SysUser) SelectByPage(queryMap map[string]any, dataScopeSQL string) ([]model.SysUser, int64) {
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
	total := int64(0)
	arr := []model.SysUser{}

	// 查询数量 长度为0直接返回
	totalSql := "select count(u.user_id) as 'total' from sys_user u " + whereSql + dataScopeSQL
	totalRows, err := db.RawDB("", totalSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
		return arr, total
	}
	total = parse.Number(totalRows[0]["total"])
	if total <= 0 {
		return arr, total
	}

	// 分页
	pageNum, pageSize := db.PageNumSize(queryMap["pageNum"], queryMap["pageSize"])
	pageSql := " limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := r.sql + whereSql + dataScopeSQL + pageSql
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
func (r SysUser) Select(sysUser model.SysUser, dataScopeSQL string) []model.SysUser {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysUser.UserId != "" {
		conditions = append(conditions, "u.user_id = ?")
		params = append(params, sysUser.UserId)
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
	querySql := r.sql + whereSql + dataScopeSQL
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysUser{}
	}

	// 转换实体
	arr := []model.SysUser{}
	if err := db.Unmarshal(rows, &arr); err != nil {
		logger.Errorf("unmarshal err => %v", err)
	}
	return arr
}

// SelectByIds 通过ID查询信息
func (r SysUser) SelectByIds(userIds []string) []model.SysUser {
	placeholder := db.KeyPlaceholderByQuery(len(userIds))
	querySql := r.sql + " where u.del_flag = '0' and u.user_id in (" + placeholder + ")"
	parameters := db.ConvertIdsSlice(userIds)
	rows, err := db.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysUser{}
	}
	// 转换实体
	arr := []model.SysUser{}
	if err := db.Unmarshal(rows, &arr); err != nil {
		logger.Errorf("unmarshal err => %v", err)
	}
	return arr
}

// Insert 新增信息
func (r SysUser) Insert(sysUser model.SysUser) string {
	// 参数拼接
	params := make(map[string]any)
	if sysUser.UserId != "" {
		params["user_id"] = sysUser.UserId
	}
	if sysUser.DeptId != "" {
		params["dept_id"] = sysUser.DeptId
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
func (r SysUser) Update(sysUser model.SysUser) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysUser.DeptId != "" {
		params["dept_id"] = sysUser.DeptId
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
	params["remark"] = sysUser.Remark
	if sysUser.UpdateBy != "" {
		params["update_by"] = sysUser.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}
	if sysUser.LoginIp != "" {
		params["login_ip"] = sysUser.LoginIp
	}
	if sysUser.LoginDate > 0 {
		params["login_date"] = sysUser.LoginDate
	}

	// 构建执行语句
	keys, values := db.KeyValueByUpdate(params)
	sql := fmt.Sprintf("update sys_user set %s where user_id = ?", keys)

	// 执行更新
	values = append(values, sysUser.UserId)
	rows, err := db.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteByIds 批量删除信息
func (r SysUser) DeleteByIds(userIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(userIds))
	username := "CASE WHEN user_name = '' THEN user_name WHEN LENGTH(user_name) >= 36 THEN CONCAT('del_', SUBSTRING(user_name, 5, 36)) ELSE CONCAT('del_', user_name) END"
	email := "CASE WHEN email = '' THEN email WHEN LENGTH(email) >= 64 THEN CONCAT('del_', SUBSTRING(email, 5, 64)) ELSE CONCAT('del_', email) END"
	phonenumber := "CASE WHEN phonenumber = '' THEN phonenumber WHEN LENGTH(phonenumber) >= 16 THEN CONCAT('del_', SUBSTRING(phonenumber, 5, 16)) ELSE CONCAT('del_', phonenumber) END"
	sql := fmt.Sprintf("update sys_user set del_flag = '1', user_name = %s, email = %s, phonenumber = %s where user_id in (%s)", username, email, phonenumber, placeholder)
	parameters := db.ConvertIdsSlice(userIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("update err => %v", err)
		return 0
	}
	return results
}

// CheckUnique 检查信息是否唯一
func (r SysUser) CheckUnique(sysUser model.SysUser) string {
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
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
		return ""
	}
	if len(rows) > 0 {
		return fmt.Sprint(rows[0]["str"])
	}
	return ""
}

// SelectByUserName 通过登录账号查询信息
func (r SysUser) SelectByUserName(userName string) model.SysUser {
	querySql := r.sql + " where u.del_flag = '0' and u.user_name = ?"
	rows, err := db.RawDB("", querySql, []any{userName})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return model.SysUser{}
	}
	// 转换实体
	arr := []model.SysUser{}
	if err := db.Unmarshal(rows, &arr); err != nil {
		logger.Errorf("unmarshal err => %v", err)
	}
	if len(arr) > 0 {
		return arr[0]
	}
	return model.SysUser{}
}

// SelectAllocatedByPage 分页查询集合By分配用户角色
func (r SysUser) SelectAllocatedByPage(query map[string]any, dataScopeSQL string) ([]model.SysUser, int64) {
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
	// 分配角色的用户
	if allocated, ok := query["allocated"]; ok && allocated != "" {
		if roleId, ok := query["roleId"]; ok && roleId != "" {
			if parse.Boolean(allocated) {
				conditions = append(conditions, `u.user_id in (
					select distinct u.user_id from sys_user u 
					inner join sys_user_role ur on u.user_id = ur.user_id 
					and ur.role_id = ?
				)`)
			} else {
				conditions = append(conditions, `u.user_id not in (
					select distinct u.user_id from sys_user u 
					inner join sys_user_role ur on u.user_id = ur.user_id 
					and ur.role_id = ?
				)`)
			}
			params = append(params, roleId)
		}
	}

	// 构建查询条件语句
	whereSql := " where u.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询结果
	total := int64(0)
	arr := []model.SysUser{}

	// 查询数量 长度为0直接返回
	totalSql := `select count(u.user_id) as 'total' from sys_user u `
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
	pageSql := " limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := r.sql + whereSql + dataScopeSQL + pageSql
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
