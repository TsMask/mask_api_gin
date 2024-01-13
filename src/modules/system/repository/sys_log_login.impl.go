package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"

	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// 实例化数据层 SysLogLoginImpl 结构体
var NewSysLogLoginImpl = &SysLogLoginImpl{
	selectSql: `select login_id, user_name, ipaddr, login_location, 
	browser, os, status, msg, login_time from sys_log_login`,

	resultMap: map[string]string{
		"login_id":       "LoginID",
		"user_name":      "UserName",
		"status":         "Status",
		"ipaddr":         "IPAddr",
		"login_location": "LoginLocation",
		"browser":        "Browser",
		"os":             "OS",
		"msg":            "Msg",
		"login_time":     "LoginTime",
	},
}

// SysLogLoginImpl 系统登录访问表 数据层处理
type SysLogLoginImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysLogLoginImpl) convertResultRows(rows []map[string]any) []model.SysLogLogin {
	arr := make([]model.SysLogLogin, 0)
	for _, row := range rows {
		SysLogLogin := model.SysLogLogin{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				datasource.SetFieldValue(&SysLogLogin, keyMapper, value)
			}
		}
		arr = append(arr, SysLogLogin)
	}
	return arr
}

// SelectSysLogLoginPage 分页查询系统登录日志集合
func (r *SysLogLoginImpl) SelectSysLogLoginPage(query map[string]any) map[string]any {
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
	result := map[string]any{
		"total": 0,
		"rows":  []model.SysLogLogin{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_log_login"
	totalRows, err := datasource.RawDB("", totalSql+whereSql, params)
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
	pageNum, pageSize := datasource.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " order by login_id desc limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := r.selectSql + whereSql + pageSql
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return result
	}

	// 转换实体
	result["rows"] = r.convertResultRows(results)
	return result
}

// SelectSysLogLoginList 查询系统登录日志集合
func (r *SysLogLoginImpl) SelectSysLogLoginList(SysLogLogin model.SysLogLogin) []model.SysLogLogin {
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
	querySql := r.selectSql + whereSql
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysLogLogin{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// InsertSysLogLogin 新增系统登录日志
func (r *SysLogLoginImpl) InsertSysLogLogin(SysLogLogin model.SysLogLogin) string {
	// 参数拼接
	params := make(map[string]any)
	params["login_time"] = time.Now().UnixMilli()
	if SysLogLogin.UserName != "" {
		params["user_name"] = SysLogLogin.UserName
	}
	if SysLogLogin.Status != "" {
		params["status"] = SysLogLogin.Status
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
	keys, values, placeholder := datasource.KeyValuePlaceholderByInsert(params)
	sql := "insert into sys_log_login (" + keys + ")values(" + placeholder + ")"

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

// DeleteSysLogLoginByIds 批量删除系统登录日志
func (r *SysLogLoginImpl) DeleteSysLogLoginByIds(loginIds []string) int64 {
	placeholder := datasource.KeyPlaceholderByQuery(len(loginIds))
	sql := "delete from sys_log_login where login_id in (" + placeholder + ")"
	parameters := datasource.ConvertIdsSlice(loginIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// CleanSysLogLogin 清空系统登录日志
func (r *SysLogLoginImpl) CleanSysLogLogin() error {
	sql := "truncate table sys_log_login"
	_, err := datasource.ExecDB("", sql, []any{})
	return err
}
