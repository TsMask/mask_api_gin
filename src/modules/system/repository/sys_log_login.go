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

// NewSysLogLogin 实例化数据层
var NewSysLogLogin = &SysLogLogin{
	sql: `select 
	login_id, user_name, ipaddr, login_location, 
	browser, os, status, msg, login_time 
	from sys_log_login`,
}

// SysLogLoginRepository 系统登录访问表 数据层处理
type SysLogLogin struct {
	sql string // 查询视图对象SQL
}

// SelectByPage 分页查询集合
func (r SysLogLogin) SelectByPage(query map[string]any) ([]model.SysLogLogin, int64) {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["ipaddr"]; ok && v != "" {
		conditions = append(conditions, "ipaddr like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["userName"]; ok && v != "" {
		conditions = append(conditions, "user_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok && v != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, v)
	}
	beginTime, ok := query["beginTime"]
	if !ok {
		beginTime, ok = query["params[beginTime]"]
	}
	if ok && beginTime != "" {
		conditions = append(conditions, "login_time >= ?")
		beginDate := date.ParseStrToDate(beginTime.(string), date.YYYY_MM_DD)
		params = append(params, beginDate.UnixMilli())
	}
	endTime, ok := query["endTime"]
	if !ok {
		endTime, ok = query["params[endTime]"]
	}
	if ok && endTime != "" {
		conditions = append(conditions, "login_time <= ?")
		endDate := date.ParseStrToDate(endTime.(string), date.YYYY_MM_DD)
		params = append(params, endDate.UnixMilli())
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询结果
	total := int64(0)
	arr := []model.SysLogLogin{}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_log_login"
	totalRows, err := db.RawDB("", totalSql+whereSql, params)
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
	pageSql := " order by login_id desc limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := r.sql + whereSql + pageSql
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
func (r SysLogLogin) Select(SysLogLogin model.SysLogLogin) []model.SysLogLogin {
	// 查询条件拼接
	var conditions []string
	var params []any
	if SysLogLogin.IPAddr != "" {
		conditions = append(conditions, "title like concat(?, '%')")
		params = append(params, SysLogLogin.IPAddr)
	}
	if SysLogLogin.UserName != "" {
		conditions = append(conditions, "user_name like concat(?, '%')")
		params = append(params, SysLogLogin.UserName)
	}
	if SysLogLogin.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, SysLogLogin.Status)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数据
	querySql := r.sql + whereSql
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysLogLogin{}
	}

	// 转换实体
	arr := []model.SysLogLogin{}
	if err := db.Unmarshal(rows, &arr); err != nil {
		logger.Errorf("unmarshal err => %v", err)
	}
	return arr
}

// Insert 新增信息
func (r SysLogLogin) Insert(SysLogLogin model.SysLogLogin) string {
	// 参数拼接
	params := make(map[string]any)
	params["login_time"] = time.Now().UnixMilli()
	if SysLogLogin.UserName != "" {
		params["user_name"] = SysLogLogin.UserName
	}
	if SysLogLogin.Status != "" {
		params["status"] = parse.Number(SysLogLogin.Status)
	}
	if SysLogLogin.IPAddr != "" {
		params["ipaddr"] = SysLogLogin.IPAddr
	}
	if SysLogLogin.LoginLocation != "" {
		params["login_location"] = SysLogLogin.LoginLocation
	}
	if SysLogLogin.Browser != "" {
		params["browser"] = SysLogLogin.Browser
	}
	if SysLogLogin.OS != "" {
		params["os"] = SysLogLogin.OS
	}
	if SysLogLogin.Msg != "" {
		params["msg"] = SysLogLogin.Msg
	}

	// 构建执行语句
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_log_login (%s)values(%s)", keys, placeholder)

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
func (r SysLogLogin) DeleteByIds(loginIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(loginIds))
	sql := fmt.Sprintf("delete from sys_log_login where login_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(loginIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// Clean 清空信息
func (r SysLogLogin) Clean() error {
	sql := "truncate table sys_log_login"
	_, err := db.ExecDB("", sql, []any{})
	return err
}
